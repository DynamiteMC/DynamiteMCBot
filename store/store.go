package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

var files = make(map[string]map[string]struct{})

func Point[T any](d T) *T {
	return &d
}

func writeFile(name string, data []byte) {
	psp := strings.Split(name, "/")
	if len(psp) > 1 {
		os.MkdirAll(strings.Join(psp[:len(psp)-1], "/"), 0755)
	}
	os.WriteFile(name, data, 0755)
}

func read(path string) map[string]struct{} {
	var output map[string]struct{}
	d, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return map[string]struct{}{}
	}
	json.Unmarshal(d, &output)
	return output
}

func addEntryTo(path string, data string) {
	if files[path] == nil {
		files[path] = read(path)
	}
	files[path][data] = struct{}{}
	d, _ := json.Marshal(files[path])
	writeFile(path, d)
}

func removeEntryFrom(path string, data string) {
	if files[path] == nil {
		files[path] = read(path)
	}
	delete(files[path], data)
	d, _ := json.Marshal(files[path])
	writeFile(path, d)
}

func isExist(path string, data string) bool {
	if files[path] == nil {
		files[path] = read(path)
	}
	if _, ok := files[path][data]; ok {
		return true
	}
	return false
}

func AddCornered(id int64) {
	addEntryTo("data/cornered.json", fmt.Sprint(id))
}

func RemoveCornered(id int64) {
	removeEntryFrom("data/cornered.json", fmt.Sprint(id))
}

func IsCornered(id int64) bool {
	return isExist("data/cornered.json", fmt.Sprint(id))
}

func AddMuted(id int64) {
	addEntryTo("data/muted.json", fmt.Sprint(id))
}

func RemoveMuted(id int64) {
	removeEntryFrom("data/muted.json", fmt.Sprint(id))
}

func IsMuted(id int64) bool {
	return isExist("data/muted.json", fmt.Sprint(id))
}
