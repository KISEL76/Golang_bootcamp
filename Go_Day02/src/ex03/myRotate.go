package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func createTarGz(srcFile, dstFile string) error {
	outFile, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	file, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = srcFile

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	return err
}

func archiveLog(filePath, archiveDir string, wg *sync.WaitGroup) {
	defer wg.Done()

	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("Ошибка чтения %s: %v\n", filePath, err)
		return
	}
	modTime := info.ModTime().Unix()

	filename := filepath.Base(filePath)
	archiveName := fmt.Sprintf("%s_%d.tar.gz", filename, modTime)

	var archivePath string
	if archiveDir != "" {
		archivePath = filepath.Join(archiveDir, archiveName)
	} else {
		archivePath = filepath.Join(filepath.Dir(filePath), archiveName)
	}

	if err := createTarGz(filePath, archivePath); err != nil {
		fmt.Printf("Ошибка архивации %s: %v\n", filePath, err)
	} else if err == nil {
		os.Remove(filePath)
	}
}

func main() {
	archiveDir := flag.String("a", "", "")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("Использование: ./myRotate -a <архивная_директория> <log1> <log2> ... <log n>")
		return
	}

	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go archiveLog(file, *archiveDir, &wg)
	}
	wg.Wait()
}
