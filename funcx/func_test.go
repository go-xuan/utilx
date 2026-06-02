package funcx

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestExecute(t *testing.T) {
	ctx := t.Context()
	ExecuteV(testV1, testV2)

	mergeX := MergeX(testX1, testX2, testX3)
	if err := mergeX(ctx); err != nil {
		fmt.Println(err)
	}

	if err := ExecuteX(ctx, testX1, testX2, testX3); err != nil {
		fmt.Println(err)
	}
}

func testV1() {
	println("this is v1")
}

func testV2() {
	println("this is v2")
}

func testC1(context.Context) {
	println("this is c1")
}

func testE1() error {
	println("this is e1")
	return nil
}

func testX1(context.Context) error {
	println("this is x1")
	return nil
}

func testX2(context.Context) error {
	println("this is x2")
	return errors.New("x2 error")
}

func testX3(context.Context) error {
	println("this is X3")
	return errors.New("x3 error")
}
