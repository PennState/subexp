package subexp

import (
	"fmt"
	"reflect"
)

const indexErrorFmt = "index out of bounds for array with len %d, index was %d"

type IndexError struct {
	Len int
	Idx int
}

func boundsCheck(v interface{}, i int) error {
	len := reflect.TypeOf(v).Len() // panics if v is not an array

	if i < 0 || i > len-1 {
		return IndexError{
			Len: len,
			Idx: i,
		}
	}

	return nil
}

func (e IndexError) Error() string {
	return fmt.Sprintf(indexErrorFmt, e.Len, e.Idx)
}

const keyErrorFmt = "requested named capture group was not found: %s"

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

func (e KeyError) Error() string {
	return fmt.Sprintf(keyErrorFmt, e.Key)
}

const noTextErrorFmt = "no captured text for name: %s"

type NoTextError struct {
	Key string
}

func textCheck(v []string, k string) error {
	if len(v) < 1 {
		return NoTextError{
			Key: k,
		}
	}

	return nil
}

func (e NoTextError) Error() string {
	return fmt.Sprintf(noTextErrorFmt, e.Key)
}
