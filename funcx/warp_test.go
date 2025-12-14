package funcx

import (
	"context"
	"fmt"
	"testing"
)

func TestWarp(t *testing.T) {
	VWarp(testV1, VDuration)()

	fmt.Print("\n\n\n")

	CWarp(testC1, CDuration)(t.Context())

	fmt.Print("\n\n\n")

	if err := EWarp(testE1, EDuration)(); err != nil {
		t.Errorf("EWarp failed: %v", err)
	}

	fmt.Print("\n\n\n")

	if err := XWarp(testX1, beforeX, afterX, XDuration)(t.Context()); err != nil {
		t.Errorf("XWarp failed: %v", err)
	}

	fmt.Print("\n\n\n")

	if err := XWarp(testX1, XDuration, afterX, beforeX)(t.Context()); err != nil {
		t.Errorf("XWarp failed: %v", err)
	}
}

func beforeX(function X) X {
	return func(ctx context.Context) error {
		fmt.Println("=====before X=====")
		if err := function(ctx); err != nil {
			return err
		}
		return nil
	}
}

func afterX(function X) X {
	return func(ctx context.Context) error {
		if err := function(ctx); err != nil {
			return err
		}
		fmt.Println("=====after X=====")
		return nil
	}
}
