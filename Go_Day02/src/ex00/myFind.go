package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func linksPrint(path, cleanPath string) {
	target, err := os.Readlink(path)
	if err != nil {
		fmt.Printf("%s -> [broken]\n", cleanPath)
	} else {
		fmt.Printf("%s -> %s\n", cleanPath, target)
	}
}

func filePrint(key int, extFlag *string, path string, info os.FileInfo) (err error) {
	if *extFlag != "" && key&2 != 0 && strings.HasSuffix(info.Name(), "."+*extFlag) {
		fmt.Println(path)
		return nil
	} else if *extFlag == "" {
		fmt.Println(path)
		return nil
	}
	return nil
}

func btoi(val *bool) int {
	if *val {
		return 1
	}
	return 0
}

func utilityPrint(key int, info os.FileInfo, path, cleanPath string, extFlag *string) (err error) {
	switch key {
	case 0: // 000
		if info.Mode()&os.ModeSymlink != 0 {
			linksPrint(path, cleanPath)
			return nil
		}
		filePrint(key, extFlag, path, info)
	case 1: // 001
		if info.Mode()&os.ModeSymlink != 0 {
			linksPrint(path, path)
			return nil
		}
	case 2: // 010
		if !info.IsDir() {
			filePrint(key, extFlag, path, info)
		}
	case 3: // 011
		if !info.IsDir() {
			filePrint(key, extFlag, path, info)
		} else if info.Mode()&os.ModeSymlink != 0 {
			linksPrint(path, path)
			return nil
		}
	case 4: // 100
		if info.IsDir() {
			fmt.Println(path)
			return nil
		}
	case 5: // 101
		if info.IsDir() {
			fmt.Println(path)
			return nil
		} else if info.Mode()&os.ModeSymlink != 0 {
			linksPrint(path, path)
			return nil
		}
	case 6: // 110
		if !info.IsDir() {
			filePrint(key, extFlag, path, info)
		} else {
			fmt.Println(path)
			return nil
		}
	case 7: // 111
		if info.Mode()&os.ModeSymlink != 0 {
			linksPrint(path, path)
			return nil
		}

		if !info.IsDir() {
			filePrint(key, extFlag, path, info)
		} else {
			fmt.Println(path)
			return nil
		}
	}
	return nil
}

func findFiles(root string, dirFlag, fileFlag, linkFlag *bool, extFlag *string) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка доступа к %s: %v\n", path, err)
			return nil
		}

		key := (btoi(dirFlag) << 2) | (btoi(fileFlag) << 1) | btoi(linkFlag)

		utilityPrint(key, info, path, root, extFlag)
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка обхода: %v\n", err)
	}
}

func main() {
	dirFlag := flag.Bool("d", false, "on/off to show directories\n")
	fileFlag := flag.Bool("f", false, "on/off to show files\n")
	linkFlag := flag.Bool("sl", false, "on/off to show symlinks\n")
	extFlag := flag.String("ext", "", "works in pair with flag f: ./myFind -f -ext '.ext' <path>\n")
	flag.Parse()

	if len(flag.Args()) == 1 {
		homeDir, _ := os.UserHomeDir()
		root, err := filepath.Abs(filepath.Join(homeDir, flag.Args()[0]))
		if err != nil {
			fmt.Printf("Ошибка нахождения пути %s: %v\n", root, err)
			return
		}

		findFiles(root, dirFlag, fileFlag, linkFlag, extFlag)
	} else {
		fmt.Printf("Неправильное использование утилиты: ./myFind <flags> <one_path>\n")
		return
	}
}
