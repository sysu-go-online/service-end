package model

import (
	"io/ioutil"
	"path/filepath"
)

var id = 1

func dfs(path string, depth int) ([]FileStructure, error) {
	var structure []FileStructure
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		tmp := FileStructure{
			ID:         id,
			Name:       file.Name(),
			Root:       false,
			IsSelected: false,
		}
		id++
		if file.IsDir() {
			tmp.Type = "dir"
			nextPath := filepath.Join(path, file.Name())
			tmp.Children, err = dfs(nextPath, depth+1)
			if err != nil {
				return nil, err
			}
		} else {
			tmp.Type = "file"
		}
		structure = append(structure, tmp)
	}
	return structure, nil
}
