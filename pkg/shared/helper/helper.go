package helper

import (
	"fmt"
	"reflect"
)

func IsZeroValue(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func FailToError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s %s\n", msg, err)
	}
}
