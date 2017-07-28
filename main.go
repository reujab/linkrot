package main

import (
	"bufio"
	"os"
	"path/filepath"
)

func main() {
	die(filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		// skip version control
		if info.Name() == ".git" {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		die(err)
		defer func() { die(file.Close()) }()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			// hyperlinks cannot be less than 11 characters
			// example of a short hyperlink: http://j.tl
			if len(line) < 11 {
				continue
			}
		}

		return nil
	}))
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}
