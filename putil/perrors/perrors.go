package perrors

import (
	"bytes"
	"errors"
	. "putil/perror"
	"strconv"
)

type implPerror struct {
	code int
	err  error
}

func (perr *implPerror) Code() int {
	return perr.code
}
func (perr *implPerror) Err() error {
	return perr.err
}
func (perr *implPerror) SetCode(code int) {
	perr.code = code
}
func (perr *implPerror) SetErr(err error) {
	perr.err = err
}
func (perr *implPerror) Error() string {
	var str bytes.Buffer
	str.WriteString("code=")
	str.WriteString(strconv.Itoa(perr.code))
	str.WriteString(" decrip:")
	str.WriteString(perr.err.Error())

	return str.String()
}
func New(code int, desp string) Perror {
	perr := new(implPerror)
	perr.err = errors.New(desp)
	perr.code = code
	return perr
}

func New2(err error) Perror {
	perr := new(implPerror)
	perr.err = err
	perr.code = -1
	return perr
}
