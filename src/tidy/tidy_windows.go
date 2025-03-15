package tidy

/*
#cgo LDFLAGS: -L../ -ltidy
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#inlcude "headers/window_call.c"
*/
import "C"
import (
	// "fmt"
	// "io"
	// "os"
	"unsafe"
)

func TidyHTML(input []byte) []byte {
	cInput := C.CBytes(input)
	defer C.free(unsafe.Pointer(cInput))

	cOutput := C.tidy_html((*C.char)(cInput), C.uint64_t(len(input)))

	if cOutput == nil {
		return []byte("Error processing HTML")
	}
	defer C.free(unsafe.Pointer(cOutput))

	return C.GoBytes(unsafe.Pointer(cOutput), C.int(C.strlen(cOutput)))
}
