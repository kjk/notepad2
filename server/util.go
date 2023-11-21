package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/kjk/common/u"
)

var (
	f          = fmt.Sprintf
	e          = fmt.Errorf
	must       = u.Must
	panicIf    = u.PanicIf
	panicIfErr = u.PanicIfErr
	isWinOrMac = u.IsWinOrMac
	isLinux    = u.IsLinux
	formatSize = u.FormatSize
)

func ctx() context.Context {
	return context.Background()
}

func runLoggedInDirMust(dir string, exe string, args ...string) {
	cmd := exec.Command(exe, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	must(err)
}

func getCallstackFrames(skip int) []string {
	var callers [32]uintptr
	n := runtime.Callers(skip+1, callers[:])
	frames := runtime.CallersFrames(callers[:n])
	var cs []string
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		s := frame.File + ":" + strconv.Itoa(frame.Line)
		cs = append(cs, s)
	}
	return cs
}

func getCallstack(skip int) string {
	frames := getCallstackFrames(skip + 1)
	return strings.Join(frames, "\n")
}

const (
	htmlMimeType     = "text/html; charset=utf-8"
	jsMimeType       = "text/javascript; charset=utf-8"
	markdownMimeType = "text/markdown; charset=UTF-8"
)

func jsonUnmarshalReader(r io.Reader, v interface{}) error {
	d, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, v)
}

func fmtSmart(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

func startLoggedInDir(dir string, exe string, args ...string) (func(), error) {
	cmd := exec.Command(exe, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return func() {
		cmd.Process.Kill()
	}, nil
}

func push[S ~[]E, E any](s *S, els ...E) {
	*s = append(*s, els...)
}

func sliceLimit[S ~[]E, E any](s S, max int) S {
	if len(s) > max {
		return s[:max]
	}
	return s
}
