package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserType interface {
	IsProjectManager() bool
	IsContributor() bool
}

type LoginUser struct {
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type TeamMember struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type Supervisor struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

type KeyQuest struct {
	Name             string `bson:"name" json:"name"`
	CompletedTasks   int    `bson:"completedTasks" json:"completedTasks"`
	IncompletedTasks int    `bson:"incompletedTasks" json:"incompletedTasks"`
	Rank             string `bson:"rank" json:"rank"`
	Exp              int    `bson:"exp" json:"exp"`
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Icon              string             `bson:"icon,omitempty" json:"icon,omitempty"`
	Wallet            int                `bson:"wallet" json:"wallet"`
	Email             string             `bson:"email" json:"email"`
	Password          string             `bson:"password" json:"password"`
	Level             int                `bson:"level" json:"level"`
	TotalExp          int                `bson:"totalExp" json:"totalExp"`
	TotalPoints       int                `bson:"totalPoints" json:"totalPoints"`
	TotalExpThisMonth int                `bson:"totalExpThisMonth" json:"totalExpThisMonth"`
	WeekStreak        int                `bson:"weekStreak" json:"weekStreak"`
	Progress          int                `bson:"progress" json:"progress"`
	Following         []string           `bson:"following" json:"following"`
	Followers         []string           `bson:"followers" json:"followers"`
	KeyQuest          KeyQuest           `bson:"keyQuest" json:"keyQuest,omitempty"`
	KeyTasks          StatusColumn       `bson:"keyTasks" json:"keyTasks"`
	Title             string             `bson:"title" json:"title"`
	Team              []TeamMember       `bson:"team" json:"team"`
	Supervisor        Supervisor         `bson:"supervisor" json:"supervisor"`
}

func (u *User) IsProjectManager() bool {
	return u.Title == "Manager (PM)"
}

func (u *User) IsContributor() bool {
	return u.Title == "Contributor (IC)"
}
