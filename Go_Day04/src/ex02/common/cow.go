package common

/*
#include <stdlib.h>
char *ask_cow(char* phrase);
*/
import "C"
import (
	"unsafe"
)

func Cowify(phrase string) string {
	cStr := C.CString(phrase)
	defer C.free(unsafe.Pointer(cStr))

	result := C.ask_cow(cStr)
	defer C.free(unsafe.Pointer(result))

	return C.GoString(result)
}
