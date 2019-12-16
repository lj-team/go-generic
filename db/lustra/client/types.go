package client

type Client interface {
	Exec(command ...string) (string, error)
	Async(command ...string) error
	Close()
}
