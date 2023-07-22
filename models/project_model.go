package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID       string `bson:"id" json:"id"`
	Name     string `bson:"name" json:"name"`
	Content  string `bson:"content" json:"content"`
	ExpYield int    `bson:"expYield" json:"expYield"`
	Status   string `bson:"status" json:"status"`
	SubTask  []Task `bson:"subTasks" json:"subTasks"`
}

type StatusColumn struct {
	Name  string `bson:"name" json:"name"`
	Items []Task `bson:"items" json:"items"`
}

type Columns struct {
	NotStarted StatusColumn `bson:"notStarted,omitempty" json:"notStarted"`
	InProgress StatusColumn `bson:"inProgress,omitempty" json:"inProgress"`
	InReview   StatusColumn `bson:"inReview,omitempty"   json:"inReview"`
	Complete   StatusColumn `bson:"complete,omitempty"   json:"complete"`
}

type Project struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Tasks            []Task             `bson:"tasks" json:"tasks"`
	CompletedTasks   int                `bson:"completedTasks" json:"completedTasks"`
	IncompletedTasks int                `bson:"incompletedTasks" json:"incompletedTasks"`
	ProjectManager   primitive.ObjectID `bson:"projectManager,omitEmpty" json:"projectManager"`
	Rank             string             `bson:"rank" json:"rank"`
	Exp              int                `bson:"exp" json:"exp"`
	TakenBy          string             `bson:"takenBy,omitempty" json:"takenBy"`
}
