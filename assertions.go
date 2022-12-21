package assert

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"text/tabwriter"
)

func NotNil(tb testing.TB, right interface{}, msgAndArgs ...interface{}) {
	if tb != nil {
		tb.Helper()
	}
	if !IsNil(right) {
		return
	}
	assertLog(tb, nil, right, false, msgAndArgs...)
}

func Nil(tb testing.TB, actual interface{}, msgAndArgs ...interface{}) {
	if tb != nil {
		tb.Helper()
	}
	if IsNil(actual) {
		return
	}
	assertLog(tb, nil, actual, true, msgAndArgs...)
}

func NotEqual(tb testing.TB, left, right interface{}, msgAndArgs ...interface{}) {
	if tb != nil {
		tb.Helper()
	}
	if !reflect.DeepEqual(left, right) {
		return
	}
	assertLog(tb, left, right, false, msgAndArgs...)
}

// Equal checks if values are equal
// Ref: gofiber/utils
func Equal(tb testing.TB, expected, actual interface{}, msgAndArgs ...interface{}) {
	if tb != nil {
		tb.Helper()
	}
	if reflect.DeepEqual(expected, actual) {
		return
	}
	assertLog(tb, expected, actual, true, msgAndArgs...)
}

func assertLog(tb testing.TB, a, b interface{}, isEqual bool, msgAndArgs ...interface{}) {
	aType := "<nil>"
	bType := "<nil>"

	if a != nil {
		aType = reflect.TypeOf(a).String()
	}
	if b != nil {
		bType = reflect.TypeOf(b).String()
	}

	testName := "Equal"
	leftTitle := "Expected"
	rightTitle := "Actual"
	if !isEqual {
		testName = "NotEqual"
		leftTitle = "Left"
		rightTitle = "Right"
	}
	if tb != nil {
		testName = fmt.Sprintf("%s(%s)", tb.Name(), testName)
	}

	_, file, line, _ := runtime.Caller(2)

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 5, ' ', 0)
	_, _ = fmt.Fprintf(w, "\nTest:\t%s", testName)
	_, _ = fmt.Fprintf(w, "\nTrace:\t%s:%d", filepath.Base(file), line)
	message := messageFromMsgAndArgs(msgAndArgs...)
	if message != "" {
		_, _ = fmt.Fprintf(w, "\nDescription:\t%s", message)
	}
	_, _ = fmt.Fprintf(w, "\n%s:\t%v\t(%s)", leftTitle, a, aType)
	_, _ = fmt.Fprintf(w, "\n%s:\t%v\t(%s)", rightTitle, b, bType)

	result := ""
	if err := w.Flush(); err != nil {
		result = err.Error()
	} else {
		result = buf.String()
	}

	if tb != nil {
		tb.Fatal(result)
	} else {
		log.Fatal(result)
	}
}

func Panics(t *testing.T, title string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			if t != nil {
				t.Fatalf("%s: didn't panic as expected", title)
			} else {
				log.Fatalf("%s: didn't panic as expected", title)
			}
		}
	}()
	f()
}

// IsNil Ref: stretchr/testify
func IsNil(o interface{}) bool {
	if o == nil {
		return true
	}

	value := reflect.ValueOf(o)
	kind := value.Kind()
	isNilableKind := containsKind(
		[]reflect.Kind{
			reflect.Chan, reflect.Func,
			reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Slice, reflect.UnsafePointer},
		kind)

	if isNilableKind && value.IsNil() {
		return true
	}

	return false
}

// containsKind checks if a specified kind in the slice of kinds.
// Ref: stretchr/testify
func containsKind(kinds []reflect.Kind, kind reflect.Kind) bool {
	for i := 0; i < len(kinds); i++ {
		if kind == kinds[i] {
			return true
		}
	}

	return false
}

// Ref: stretchr/testify
func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}