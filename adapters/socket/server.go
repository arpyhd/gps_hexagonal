package socket

import (
	logger "gps_hexagonal/helpers/logger"
	"net"
	"os"
	"syscall"
)

const (
	EPOLLET        = 1 << 31
	MaxEpollEvents = 1024
)

var (
	SERVER_HOST string
	SERVER_PORT int
)

var FdList []int

func Server() {
	var event syscall.EpollEvent
	var events [MaxEpollEvents]syscall.EpollEvent

	//Declare socket Nonblocking
	fd, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		logger.Log.Println("Server error :", err)
		os.Exit(1)
	}
	defer syscall.Close(fd)

	logger.Log.Println("Server Socket fd :", fd)

	// Set Port
	addr := syscall.SockaddrInet4{Port: SERVER_PORT}
	copy(addr.Addr[:], net.ParseIP(SERVER_HOST).To4())

	if err = syscall.SetNonblock(fd, true); err != nil {
		logger.Log.Println("Server set nonblock: ", err)
		os.Exit(1)
	}

	// Start listen on port
	syscall.Bind(fd, &addr)
	syscall.Listen(fd, syscall.SOMAXCONN)

	// Create epoll
	epfd, e := syscall.EpollCreate1(0)
	if e != nil {
		logger.Log.Println("Server error epoll_create1 :", e)
		os.Exit(1)
	}
	logger.Log.Println("Server Epoll :", epfd)
	defer syscall.Close(epfd)

	// Declare Events
	event.Events = syscall.EPOLLIN | EPOLLET | syscall.EPOLLERR | syscall.EPOLLHUP
	event.Fd = int32(fd)

	// Configure epoll to listen on events
	if e = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event); e != nil {
		logger.Log.Println("Server epoll_ctl :", e)
		os.Exit(1)
	}

	// Parse events
	for {
		logger.Log.Println("Server waiting")
		new_events, e := syscall.EpollWait(epfd, events[:], -1)
		logger.Log.Println("Server new_events", new_events)
		if e != nil {
			logger.Log.Println("Server epoll_wait  error:", e)
			os.Exit(1)
		}

		for ev := 0; ev < new_events; ev++ {
			if (events[ev].Events&syscall.EPOLLERR != 0) ||
				(events[ev].Events&syscall.EPOLLHUP != 0) ||
				(events[ev].Events&syscall.EPOLLIN) == 0 {

				/* Error condition */
				logger.Log.Println("Server error : epoll error")
				syscall.EpollCtl(epfd, syscall.EPOLL_CTL_DEL, int(events[ev].Fd), &event)
				syscall.Close(int(events[ev].Fd))
				continue

			} else if int(events[ev].Fd) == fd {
				for {
					connFd, _, err := syscall.Accept(fd)
					if err != nil {
						if (err == syscall.EAGAIN) || (err == syscall.EWOULDBLOCK) {
							logger.Log.Println("Server Processed all incoming connections")
							break
						} else {
							logger.Log.Println("Server error accept :", e)
							break
						}
					}
					err = syscall.SetNonblock(connFd, true)
					if err != nil {
						logger.Log.Println("Server error SetNonblock :", e)
						syscall.Close(int(connFd))
					}
					logger.Log.Println("Server append in FdList :", connFd)
					FdList = append(FdList, connFd)
					logger.Log.Println("Server Fdlist :", FdList)
				}
			}

		}
	}
}
