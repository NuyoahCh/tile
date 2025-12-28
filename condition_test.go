package tile

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIfThenElse(t *testing.T) {
	i := 7
	i = IfThenElse(false, i, 0)
	assert.Equal(t, i, 0)
}

func ExampleIfThenElse() {
	result := IfThenElse(true, "yes", "no")
	fmt.Println(result)

	// Output:
	// yes
}

func TestIfThenElseFunc(t *testing.T) {
	resp, err := IfThenElseFunc(true, func() (int, error) {
		return 0, nil
	}, func() (int, error) {
		return 1, errors.New("some error")
	})
	assert.NoError(t, err)
	assert.Equal(t, resp, 0)
}

func ExampleIfThenElseFunc() {
	code, err := IfThenElseFunc(false, func() (code int, err error) {
		// do something when condition is true
		// ...
		return 0, nil
	}, func() (code int, err error) {
		// do something when condition is false
		// ...
		return 1, errors.New("some error when execute func2")
	})
	fmt.Println(code)
	fmt.Println(err)

	// Output:
	// 1
	// some error when execute func2
}
