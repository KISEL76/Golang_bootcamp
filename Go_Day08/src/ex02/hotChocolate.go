package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -framework Cocoa
#include "application.h"
#include "window.h"
*/
import "C"
import "unsafe"

func main() {
	// Инициализация приложения
	C.InitApplication()

	// Создание окна 300x200
	title := C.CString("School 21")
	defer C.free(unsafe.Pointer(title))

	win := C.Window_Create(100, 100, 300, 200, title)

	// Показать окно
	C.Window_MakeKeyAndOrderFront(win)

	// Запустить event loop
	C.RunApplication()
}
