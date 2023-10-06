package ports

type Gps interface {
	Read(buf []byte) error
	Send() error
}
