// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

type Student struct {
	//ID     		primitive.ObjectID 	`bson:"_id,omitempty"`
	Name       string                 `bson:"name,omitempty"`
	Age        int                    `bson:"age,omitempty"`
	Subject    string                 `bson:"subject,omitempty"`
	CreatedAt  time.Time              `bson:"createdAt,omitempty"`
	CustomInfo map[string]interface{} `bson:"customInfo,omitempty"`
}

func StudentRecordTest(c *gin.Context) {
	c.String(http.StatusOK, "StudentRecordTest!")
	collName := "student"
	_, errVal := mongoHndl.CreateIndex(collName, "Name")
	if errVal != nil {
		log.Println("Create index failed on Name field : ", errVal)
	}

	//add document to student collection.
	insertStudentInDB(collName, "Osman Amjad", 21)
	//update document in student collection.
	insertStudentInDB(collName, "Osman Amjad", 22)
	//fetch document from student db based on index
	student, err := getStudentFromDB(collName, "Osman Amjad")
	if err == nil {
		log.Printf("Retrieved student %v ", student)
	} else {
		log.Printf("Failed to retrieve student %v. Error - ", student, err)
	}

	insertStudentInDB(collName, "John Smith", 25)

	// test document fetch from student that doesn't exist.
	qName := "Nerf Doodle"
	student, err = getStudentFromDB(collName, qName)
	if err == nil {
		log.Printf("Retrieved student %v ", qName)
	} else {
		log.Printf("Failed to retrieve student %v. Error - ", qName, err)
	}
	c.JSON(http.StatusOK, gin.H{})
}

func insertStudentInDB(collName string, name string, age int) {
	student := Student{
		Name:      name,
		Age:       age,
		CreatedAt: time.Now(),
	}
	filter := bson.M{}
	_, err := mongoHndl.PutOneCustomDataStructure(collName, filter, student)
	if err != nil {
		log.Printf("Inserting student %v failed with error ", student, err)
		return
	}
	log.Printf("Inserting student %v successful ", student)
}

func getStudentFromDB(collName string, name string) (Student, error) {
	var student Student
	filter := bson.M{}
	filter["name"] = name

	result, err := mongoHndl.GetOneCustomDataStructure(collName, filter)

	if err == nil {
		bsonBytes, _ := bson.Marshal(result)
		bson.Unmarshal(bsonBytes, &student)

		return student, nil
	}
	return student, err
}
