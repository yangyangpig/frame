package perror

type Perror interface {
	Code() int
	Err() error
	SetCode(code int)
	SetErr(err error)
	Error() string
}
