package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"unicode"

	"github.com/mvdan/xurls"
	"github.com/urfave/cli"
)

func cmdWalk(ctx *cli.Context) {
	dir := "."
	if len(ctx.Args()) != 0 {
		dir = ctx.Args()[0]
	}

	excludes := make([]*regexp.Regexp, len(ctx.StringSlice("exclude")))
	for i, exclude := range ctx.StringSlice("exclude") {
		excludes[i] = regexp.MustCompile(exclude)
	}

	die(filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// skip version control
		if info.Name() == ".git" {
			return filepath.SkipDir
		}

		for _, regex := range excludes {
			if regex.MatchString(path) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		die(err)
		defer func() { die(file.Close()) }()

		var linesScanned int
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			linesScanned++
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

			for _, match := range xurls.Strict.FindAllString(line, -1) {
				uri, err := url.Parse(match)
				die(err)
				waitgroup.Add(1)
				go queue(fmt.Sprintf("%s:%d", path, linesScanned), uri)
			}
		}
		// don't check scanner.Err() because some lines might be too long

		return nil
	}))

	waitgroup.Wait()
}
