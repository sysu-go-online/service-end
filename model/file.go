package model

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileStructure defines the structure of file
type FileStructure struct {
	ID         int             `json:"id"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Children   []FileStructure `json:"children"`
	Root       bool            `json:"root"`
	IsSelected bool            `json:"isSelected"`
}

// ROOT defines the root directory
var ROOT = "/home"

// GetProjectNameByID return project according to the given id
func GetProjectNameByID(id string) (string, error) {
	// TODO:
	return "test", nil
}

// UpdateFileContent update content with given filepath and content
func UpdateFileContent(projectName string, username string, filePath string, content string, create bool, dir bool) error {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectName)
	if err != nil {
		return err
	}
	absPath := getFilePath(username, projectName, filePath)

	// Update file, if the file not exists, judge accroding to the given param
	if create {
		if dir {
			err = os.Mkdir(filePath, os.ModeAppend)
		} else {
			err = ioutil.WriteFile(absPath, []byte(content), os.ModeAppend)
		}
	} else {
		err = ioutil.WriteFile(absPath, []byte(content), os.ModeAppend)
	}
	return err
}

// DeleteFile delete file accroding to the given path
func DeleteFile(projectName string, username string, filePath string) error {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectName)
	if err != nil {
		return err
	}
	absPath := getFilePath(username, projectName, filePath)

	// Delete file
	err = os.RemoveAll(absPath)
	return err
}

// GetFileContent returns required file content
func GetFileContent(projectName string, username string, filePath string) ([]byte, error) {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectName)
	if err != nil {
		return nil, err
	}
	absPath := getFilePath(username, projectName, filePath)

	// Read file content
	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// GetFileStructure read file structure and return it
func GetFileStructure(projectName string, username string) (*FileStructure, error) {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectName)
	if err != nil {
		return nil, err
	}
	absPath := getFilePath(username, projectName, "/")

	// Recurisively get file structure
	s, err := dfs(absPath, 0)
	if err != nil {
		return nil, err
	}
	// Add root content
	root := FileStructure{
		ID:         1,
		Name:       projectName,
		Type:       "dir",
		Children:   s,
		Root:       true,
		IsSelected: true,
	}
	return &root, nil
}

func getFilePath(username string, projectName string, filePath string) string {
	return filepath.Join(ROOT, username, "go/src/github.com", projectName, filePath)
}
