package subexp

import (
	"fmt"
	"reflect"
)

const indexErrorFmt = "index out of bounds for array with len %d, index was %d"

/*
IndexError can occur when bounds checking is performed while retrieving
capture groups using the ByIndex function.
*/
type IndexError struct {
	Len int
	Idx int
}

func boundsCheck(v interface{}, i int) error {
	len := reflect.ValueOf(v).Len() // panics if v is not an array

	if i < 0 || i > len-1 {
		return IndexError{
			Len: len,
			Idx: i,
		}
	}

	return nil
}

/*
Error implements the built-in error interface.

See: https://github.com/golang/go/blob/go1.14.4/src/builtin/builtin.go#L260
*/
func (e IndexError) Error() string {
	return fmt.Sprintf(indexErrorFmt, e.Len, e.Idx)
}

const keyErrorFmt = "requested named capture group was not found: %s"

/*
KeyError can occur when trying to retrieve one or more named capture
groups using the AllByName() or FirstByName() methods.
*/
type KeyError struct {
	Key string
}

func keyCheck(v map[string][]string, k string) error {
	if _, ok := v[k]; !ok {
		return KeyError{
			Key: k,
		}
	}

	return nil
}

/*
Error implements the built-in error interface.

See: https://github.com/golang/go/blob/go1.14.4/src/builtin/builtin.go#L260
*/
func (e KeyError) Error() string {
	return fmt.Sprintf(keyErrorFmt, e.Key)
}

const noTextErrorFmt = "no captured text for name: %s"

/*
NoTextError can occur when attempting to retrieve the text matching an
optional named capture group.
*/
type NoTextError struct {
	Key string
}

func textCheck(v []string, k string) error {
	if len(v) < 1 || len(v[0]) == 0 {
		return NoTextError{
			Key: k,
		}
	}

	return nil
}

/*
Error implements the built-in error interface.

See: https://github.com/golang/go/blob/go1.14.4/src/builtin/builtin.go#L260
*/
func (e NoTextError) Error() string {
	return fmt.Sprintf(noTextErrorFmt, e.Key)
}
