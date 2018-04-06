package cmd

import (
	"strings"
	"os"
	"bytes"
	"io"
	"path/filepath"
	"text/template"
	log "github.com/sirupsen/logrus"
)
// isEmpty checks if a given path is empty.
// Hidden files in path are ignored.
func isEmpty(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	if !fi.IsDir() {
		return fi.Size() == 0
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil && err != io.EOF {
		log.Fatal(err.Error())
	}

	for _, name := range names {
		if len(name) > 0 && name[0] != '.' {
			return false
		}
	}
	return true
}
// exists checks if a file or directory exists.
func exists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		log.Fatal(err.Error())
	}
	return false
}

func executeTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{"comment": commentifyString}).Parse(tmplStr)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	return buf.String(), err
}

func writeStringToFile(path string, s string) error {
	return writeToFile(path, strings.NewReader(s))
}

// writeToFile writes r to file with path only
func writeToFile(path string, r io.Reader) error {
	dir := filepath.Dir(path)
	if dir != "" {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return err
}

// commentfyString comments every line of in.
func commentifyString(in string) string {
	var newlines []string
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "//") {
			newlines = append(newlines, line)
		} else {
			if line == "" {
				newlines = append(newlines, "//")
			} else {
				newlines = append(newlines, "// "+line)
			}
		}
	}
	return strings.Join(newlines, "\n")
}