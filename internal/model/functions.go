package model

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func Debounce(f func(), delay time.Duration) func() {
	var timer *time.Timer
	return func() {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(delay, f)
	}
}

func Debounce2[T any, K any](f func(one T, second K), delay time.Duration) func(one T, second K) {
	var timer *time.Timer
	return func(one T, second K) {
		if timer != nil {
			timer.Stop()
		}
		var arg = func() {
			f(one, second)
		}
		timer = time.AfterFunc(delay, arg)
	}
}

type DateFromFileNameParams struct {
	URI         string
	ProjectRoot string
}

func DateFromFile(params DateFromFileNameParams) (time.Time, error) {
	file, err := uriToFilePath(params.URI)
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed to parse URI: %w", err)
	}

	relativePath, err := filepath.Rel(params.ProjectRoot, file)
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed to get relative path: %w", err)
	}
	date, err := time.Parse("2006/01/02.tsf", relativePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed to parse date: %w", err)
	}

	return date, nil
}

func uriToFilePath(uri string) (string, error) {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	// Convert URL path to system-specific file path
	filePath := parsedURL.Path
	if runtime.GOOS == "windows" {
		// On Windows, remove the leading "/"
		if filepath.IsAbs(filePath) && len(filePath) > 0 && filePath[0] == '/' {
			filePath = filePath[1:]
		}
	}

	return filePath, nil
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Map[T any, V any](list []T, f func(T) V) []V {
	var result = make([]V, len(list))
	for i, item := range list {
		result[i] = f(item)
	}
	return result
}

type Token struct {
	Index int
	Word  string
}

func TokenizeFromIndex(input string, j int) []Token {
	// Ensure j is within bounds
	if j < 0 || j >= len(input) {
		return nil
	}

	substr := input[j:]

	var tokens []Token
	var temp Token = Token{}
	for i, r := range substr {
		if r == ' ' {
			if temp.Word != "" {
				tokens = append(tokens, temp)
				temp = Token{}
			}
		} else {
			if temp.Word == "" {
				temp.Index = i + j
			}
			temp.Word += string(r)
		}
	}
	if temp.Word != "" {
		tokens = append(tokens, temp)
	}

	return tokens
}
