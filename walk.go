package main

import (
	"bufio"
	"net/url"
	"os"
	"path/filepath"
	"unicode"

	"github.com/mvdan/xurls"
)

func walk() {
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

			// loop through every unicode character to see if the line is mostly binary data
			var binaryBytes int
			for _, char := range []rune(line) {
				if !unicode.IsPrint(char) && !unicode.IsSpace(char) {
					binaryBytes += len(string(char))
					if float32(binaryBytes)/float32(len(line)) >= 0.5 {
						// if a line has more binary bytes than printable bytes, the file is probably binary
						return nil
					}
				}
			}

			urls := xurls.Strict.FindAllString(line, -1)
			for _, match := range urls {
				uri, err := url.Parse(match)
				die(err)
				waitgroup.Add(1)
				go queue(uri)
			}
		}
		// don't check scanner.Err() because some lines might be too long

		return nil
	}))

	waitgroup.Wait()
}
