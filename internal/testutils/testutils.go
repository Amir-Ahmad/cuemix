package testutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/txtar"
)

// Convert an archive to a map where the keys are the filenames
func ArchiveToMap(a *txtar.Archive) map[string][]byte {
	outMap := make(map[string][]byte)
	for _, file := range a.Files {
		outMap[file.Name] = file.Data
	}
	return outMap
}

// ptrInt is a helper function to create an integer pointer literal
func PtrInt(i int) *int {
	return &i
}

// Write txtar file to a temp directory
func WriteTxtarToTmp(t testing.TB, txtarPath string) (string, func()) {
	t.Helper()

	a, err := txtar.ParseFile(txtarPath)
	if err != nil {
		t.Fatalf("txtar: failed to read testdata: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "cuemix_test")
	if err != nil {
		t.Fatalf("txtar: failed to create temporary directory: %v", err)
	}

	err = txtar.Write(a, tmpDir)
	if err != nil {
		t.Fatalf("txtar: failed to write files: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// messageToStr formats an input message as a string
func messageToStr(messageArgs ...interface{}) string {
	if len(messageArgs) == 1 {
		if messageStr, ok := messageArgs[0].(string); ok {
			return messageStr
		}
		return fmt.Sprintf("%+v", messageArgs[0])
	}
	if len(messageArgs) > 1 {
		return fmt.Sprintf(messageArgs[0].(string), messageArgs[1:]...)
	}

	return ""
}

// NotErr asserts that err is nil
func NotErr(t testing.TB, err error, messageArgs ...interface{}) {
	t.Helper()

	if err == nil {
		return
	}

	message := messageToStr(messageArgs...)

	if message == "" {
		t.Fatalf("got error: %v", err)
	} else {
		t.Fatalf("%s: %v", message, err)
	}
}

// IsErr asserts that err is non nil
func IsErr(t testing.TB, err error, messageArgs ...interface{}) {
	t.Helper()

	if err != nil {
		return
	}

	message := messageToStr(messageArgs...)

	if message == "" {
		t.Fatalf("expected an error")
	} else {
		t.Fatalf(message)
	}
}

// Equal asserts that expected = got
func Equal[T comparable](t testing.TB, expected, got T) {
	t.Helper()

	if expected != got {
		t.Fatalf("expected: %v; got: %v", expected, got)
	}
}
