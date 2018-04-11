package service

import "github.com/sysu-go-online/service-end/model/entities"

func InsertProject(project entities.Project) (int64, error){
	affect, err := entities.Engine.InsertOne(&project)
	if err != nil {
		return 0, err
	}
	return affect, err
}

// find project
func FindProjectByUserID(userid string) ([]entities.Project, error) {
	entities.Engine.Id(name)
	as := []entities.Project{}
	err := entities.Engine.Where("user_id=?", userid).Find(&as)
	return as, err
}
