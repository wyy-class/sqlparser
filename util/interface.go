package util

type DB interface {
	Write([]byte, int) (int, error)
	Read() ([]string, error)
	Close() error
}
