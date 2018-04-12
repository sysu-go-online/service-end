package service

import (
	"fmt"

	"github.com/sysu-go-online/service-end/model/entities"
)

func InsertProject(project entities.Project) (int64, error) {
	affect, err := entities.Engine.InsertOne(&project)
	if err != nil {
		return 0, err
	}
	return affect, err
}

// FindProjectByUserID : find project
func FindProjectByUserID(userid string) ([]entities.Project, error) {
	fmt.Printf("Enter?\n")
	as := []entities.Project{}
	err := entities.Engine.Where("user_id=?", userid).Find(&as)
	return as, err
}
