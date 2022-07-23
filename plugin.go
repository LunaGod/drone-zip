package main

import (
	"io"
	"os"
	"strings"

	glob "github.com/bmatcuk/doublestar/v4"
	"github.com/klauspost/compress/zip"
	"github.com/sirupsen/logrus"
)

type Plugin struct {
	Input  []string
	Output string
}

func (p Plugin) Exec() error {

	if len(p.Input) == 0 {
		logrus.Fatalf("please enter the file or directory to be packed")
	}

	if p.Output == "" {
		logrus.Fatalf("please enter the zip output path")
	}

	var (
		input []string
	)

	for _, inputPath := range p.Input {
		matchedPath, err := glob.FilepathGlob(inputPath)
		if err != nil {
			logrus.Fatalf("glob error: %v", err)
		}

		for _, path := range matchedPath {
			if !IsDir(path) {
				input = append(input, path)
			}
		}
	}

	logrus.Infof("input path: %v", input)

	Zip(p.Output, input)
	return nil
}

func Zip(fileName string, inputList []string) {
	fw, err := os.Create(fileName)

	if err != nil {
		logrus.Fatalf("create %s error: %v", fileName, err)
	}

	w := zip.NewWriter(fw)
	defer w.Close()

	for _, input := range inputList {
		targetFile, err := w.Create(input)
		if err != nil {
			logrus.Fatalf("create %s file error: %v", input, err)
		}

		sourceFile, err := os.Open(input)
		if err != nil {
			logrus.Fatalf("open %s file error: %v", input, err)
		}
		defer sourceFile.Close()

		_, err = io.Copy(targetFile, sourceFile)
		if err != nil {
			logrus.Fatalf("compression %s file error: %v", input, err)
		}

		logrus.Infof("compression %s file success", input)
	}
}

func Contains(s []string, item string) bool {
	for _, str := range s {
		if str == item || strings.Contains(str, item) {
			return true
		}
	}
	return false
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsDir()
}
