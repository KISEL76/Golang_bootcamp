package main

import (
	"errors"
	"fmt"
	"unsafe"
)

func getElement(arr []int, idx int) (int, error) {
	if len(arr) == 0 {
		return 0, errors.New("array is empty")
	}
	if idx >= len(arr) || idx < 0 {
		return 0, errors.New("index out of range")
	}

	basePtr := unsafe.Pointer(&arr[0])
	elementPtr := (*int)(unsafe.Pointer(uintptr(basePtr) + uintptr(idx)*unsafe.Sizeof(arr[0])))

	return *elementPtr, nil
}

func main() {
	arr := []int{10, 20, 30, 40, 50}
	value, err := getElement(arr, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", value)
	}
}
