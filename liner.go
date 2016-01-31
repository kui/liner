package liner

import (
	"bufio"
	"io"
	"log"
)

type LiningWriter struct {
	io.Writer
	closer io.Closer
	done   <-chan error
}

type LineProcessor func(line string) (err error)
type ErrorHandler func(err error)

var defaultLineProcessor = LineProcessor(func(line string) (err error) {
	log.Println(line)
	return
})
var defaultErrorHandler = ErrorHandler(func(err error) {
	log.Printf("error: %s", err.Error())
})

// NewLiningWriter return a `io.Writer` which contains
// a line processing callback.
//
// If lp is nil, use default line processor which just print with log.Println
// If eh is nil, use default error handler which just print with log.Println
func NewLiningWriter(lp LineProcessor, eh ErrorHandler) *LiningWriter {
	if lp == nil {
		lp = defaultLineProcessor
	}
	if eh == nil {
		eh = defaultErrorHandler
	}

	r, w := io.Pipe()
	done := make(chan error)
	lw := &LiningWriter{w, w, done}

	go func() {
		defer r.Close()

		var err error
		s := bufio.NewScanner(r)
		for s.Scan() {
			err = lp(s.Text())
			if err != nil {
				eh(err)
			}
		}

		err = s.Err()
		if err != nil {
			eh(err)
		}

		done <- err
	}()

	return lw
}

func (lw *LiningWriter) Close() (err error) {
	err = lw.closer.Close()
	<-lw.done
	return
}
