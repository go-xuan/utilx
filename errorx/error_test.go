package errorx

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	fmt.Println(newError(1))
	fmt.Println(newError(2))
	fmt.Println(newError(3))
	fmt.Println(newError(4))
	fmt.Println(newError(5))

	fmt.Print("-----------------\n\n\n")
	err := newError(2)
	fmt.Println(noWarp(noWarp(err)))
	fmt.Println(fmtWarp(fmtWarp(err)))
	fmt.Println(doWrap(doWrap(err)))

	fmt.Print("-----------------\n\n\n")
	err = newError(1)
	fmt.Println(noWarp(noWarp(err)))
	fmt.Println(fmtWarp(fmtWarp(err)))
	fmt.Println(doWrap(doWrap(err)))

}

func newError(e int) error {
	switch e {
	case 1:
		return New("携带栈信息的error")
	case 2:
		return errors.New("普通error")
	case 3:
		return fmt.Errorf("fmt.Errorf")
	case 4:
		return fmt.Errorf("fmt.Errorf ==> %w", errors.New("普通error"))
	case 5:
		return Wrap(errors.New("普通error"), "doWrap error")
	default:
		return nil
	}
}

func noWarp(err error) error {
	return err
}

func fmtWarp(err error) error {
	return fmt.Errorf("fmt Errorf %w", err)
}

func doWrap(err error) error {
	return Wrap(err, "do wrap error")
}
