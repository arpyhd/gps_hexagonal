package socket

import (
	logger "gps_hexagonal/helpers/logger"
	"gps_hexagonal/ports"
	"os"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
)

func consumeFdList(epfd int) {
	for len(FdList) > 0 {
		fd := FdList[0]
		logger.Log.Println("Worker use fd:", fd)
		FdList = FdList[1:]
		var event syscall.EpollEvent
		//event.Events = syscall.EPOLLIN | syscall.EPOLLOUT | EPOLLET
		event.Events = syscall.EPOLLIN | syscall.EPOLLOUT | EPOLLET | syscall.EPOLLERR | syscall.EPOLLHUP
		event.Fd = int32(fd)
		if e := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event); e != nil {
			logger.Log.Println("Worker[", epfd, "] epoll_ctl error: ", e)
			//os.Exit(1)
		}
	}

}

func Worker(gpm ports.GpsMap, gp ports.GpsPersistent) {
	var events [MaxEpollEvents]syscall.EpollEvent
	epfd, e := syscall.EpollCreate1(0)
	if e != nil {
		logger.Log.Println("Error worker epoll_create1: ", e)
		os.Exit(1)
	}
	logger.Log.Println("Worker epoll id: ", epfd)
	j := 0
	for {
		consumeFdList(epfd)
		new_events, e := syscall.EpollWait(epfd, events[:], 10)
		if e != nil {
			logger.Log.Println("Worker[", epfd, "] epoll_wait error: ", e)
			continue
		}
		for i := 0; i < new_events; i++ {
			logger.Log.Println("Worker[", epfd, "] new_event", events[i].Fd)

			if (events[i].Events&syscall.EPOLLERR != 0) ||
				((events[i].Events&syscall.EPOLLIN == 0) &&
					(events[i].Events&syscall.EPOLLOUT == 0)) {
				logger.Log.Println("Worker[", epfd, "] epoll error ", events[i].Fd)
				CloseFd(epfd, events[i].Fd)
				continue
			}

			if (events[i].Events&syscall.EPOLLHUP != 0) ||
				((events[i].Events&syscall.EPOLLIN == 0) &&
					(events[i].Events&syscall.EPOLLOUT == 0)) {
				logger.Log.Println("Worker[", epfd, "] epoll error HUP ", events[i].Fd)
				CloseFd(epfd, events[i].Fd)
				continue
			}
			if events[i].Events&syscall.EPOLLOUT != 0 {
				EpollOut(epfd, events[i].Fd, gpm)
			}

			if events[i].Events&syscall.EPOLLIN != 0 {
				EpollIn(epfd, events[i].Fd, gpm, gp)
			}

		}
		if j >= 2*500*60*15 {
			//Clear();
			j = 0
		}
		j++
	}
}

func CloseFd(epfd int, intFd int32) {
	logger.Log.Println("Worker[", epfd, "] close Fd :", intFd)
	if e := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_DEL, int(intFd), nil); e != nil {
		logger.Log.Println("Worker[", epfd, "] Error close [", intFd, "] epoll_ctl: ", e)
	}
	syscall.Close(int(intFd))
}

func EpollIn(epfd int, intFd int32, gpm ports.GpsMap, gp ports.GpsPersistent) {
	logger.Log.Println("Worker[", epfd, "] Fd[", intFd, "] call epoll in")
	buf := make([]byte, 1500)
	var message []byte
	for {
		nbytes, _ := syscall.Read(int(intFd), buf[:])
		if nbytes <= 0 {
			break
		}
		message = append(message, buf[:nbytes]...)
	}
	gpm.Add(int(intFd), gp)
	gps, err := gpm.Get(int(intFd))
	if err == nil {
		gps.Read(message)
	}
	logger.Log.Println("Received byte size: ", len(message))
}

func EpollOut(epfd int, intFd int32, gpm ports.GpsMap) {
	logger.Log.Println("Worker[", epfd, "] Fd[", intFd, "] call epoll out")
	gps, err := gpm.Get(int(intFd))
	if err == nil {
		gps.Send()
	}
}
