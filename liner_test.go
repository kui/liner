package liner

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func gen(lineNum int) (file string, err error) {
	file = filepath.Join(os.TempDir(), "test.txt")

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	for i := 0; i < lineNum; i++ {
		if _, err = fmt.Fprintln(f, "foo"); err != nil {
			return
		}
	}

	return
}

func TestLiningWriterShouldProcessLines(t *testing.T) {
	lineNum := 11
	f, err := gen(lineNum)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(f)

	t.Logf("Open %s\n", f)
	fp, err := os.Open(f)
	if err != nil {
		t.Error(err)
		return
	}
	defer fp.Close()

	c := 0
	lw := NewLiningWriter(
		func(line string) (err error) {
			c++
			return
		},
		func(err error) {
			t.Error(err)
		},
	)
	defer func() {
		if err := lw.Close(); err != nil {
			t.Error(err)
		}
	}()

	written, err := io.Copy(lw, fp)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Read data: %d bytes\n", written)

	if c != lineNum {
		t.Errorf("Expected to invoke the line processor in %d times", lineNum)
		return
	}
}

func TestLiningWriterShouldHandleError(t *testing.T) {
	lineNum := 11
	f, err := gen(lineNum)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(f)

	t.Logf("Open %s\n", f)
	fp, err := os.Open(f)
	if err != nil {
		t.Error(err)
		return
	}
	defer fp.Close()

	c := 0
	lw := NewLiningWriter(
		func(line string) (err error) {
			err = errors.New("stub error")
			return
		},
		func(err error) {
			c++
		},
	)
	defer func() {
		if err := lw.Close(); err != nil {
			t.Error(err)
		}
	}()

	written, err := io.Copy(lw, fp)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Read data: %d bytes\n", written)

	if c != lineNum {
		t.Errorf("Expected to invoke the line processor in %d times", lineNum)
		return
	}
}
