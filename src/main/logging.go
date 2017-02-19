package main

import (
  "bytes"
  "fmt"
  "runtime"
  "text/tabwriter"
  "log"
)

// CodePoint represents a point in the source code
type CodePoint struct {
  Function string
  File     string
  Line     int
}

// GetTrace returns a list of CodePoints of where the function is executing
func GetTrace(skip int) []CodePoint {
  callers := make([]uintptr, 10)
  runtime.Callers(skip+2, callers)
  var points []CodePoint

  for _, caller := range callers {
    if caller != 0 {
      fcn := runtime.FuncForPC(caller)
      point := CodePoint{Function: fcn.Name()}
      point.File, point.Line = fcn.FileLine(caller)
      points = append(points, point)
    }
  }

  return points
}

// GetFormattedTrace returns a string containing a trace
func GetFormattedTrace(skip int) string {
  trace := GetTrace(skip + 1)

  var buf bytes.Buffer
  w := tabwriter.NewWriter(&buf, 4, 0, 3, ' ', 0)
  fmt.Fprintln(&buf, "STACK TRACE")
  fmt.Fprintln(w, "File\tLine\tFunction")
  for _, point := range trace {
    fmt.Fprintf(w, "\n@ %s\t%d\t%s", point.File, point.Line, point.Function)
  }
  w.Flush()

  return buf.String()
}

// LogError logs an error and a formatted trace
func LogError(err error) {
  LogErrorMessage(err.Error())
}

// LogErrorMessage logs an error message and a formatted trace
func LogErrorMessage(msg string) {
  log.Printf("%s\n%s", msg, GetFormattedTrace(1))
}

// LogErrorMessageFatal logs an error and ends execution of the program
func LogErrorFatal(err error) {
  LogErrorMessageFatal(err.Error())
}

// LogErrorMessageFatal logs an error message and ends execution of the program
func LogErrorMessageFatal(msg string) {
  log.Fatalf("%s\n%s", msg, GetFormattedTrace(1))
}