package funcx

import (
	"context"
	"fmt"
	"testing"
)

func TestWarp(t *testing.T) {
	WrapV(testV1, DurationV)()

	fmt.Print("\n\n\n")

	WrapC(testC1, DurationC)(t.Context())

	fmt.Print("\n\n\n")

	if err := WrapE(testE1, DurationE)(); err != nil {
		t.Errorf("WrapE failed: %v", err)
	}

	fmt.Print("\n\n\n")

	if err := WrapX(testX1, beforeX, afterX, DurationX)(t.Context()); err != nil {
		t.Errorf("WrapX failed: %v", err)
	}

	fmt.Print("\n\n\n")

	if err := WrapX(testX1, DurationX, afterX, beforeX)(t.Context()); err != nil {
		t.Errorf("WrapX failed: %v", err)
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
