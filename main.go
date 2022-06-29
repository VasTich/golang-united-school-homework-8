package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

var operationKey string = "operation"
var itemKey string = "item"
var fileNameKey string = "fileName"
var idKey string = "id"

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func printFile(fileName string, writer io.Writer) error {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func performAddCommand(args Arguments, writer io.Writer) error {
	fileName, ok := args[fileNameKey]
	if !ok || len(fileName) == 0 {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	item, ok := args[itemKey]
	if !ok || len(item) == 0 {
		return fmt.Errorf("-item flag has to be specified")
	}

	var addedUser User
	err := json.Unmarshal([]byte(item), &addedUser)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var users []User
	if len(bytes) != 0 {
		err = json.Unmarshal(bytes, &users)
		if err != nil {
			return err
		}
	}

	for _, v := range users {
		if v.Id == addedUser.Id {
			writer.Write([]byte(fmt.Sprintf("Item with id %s already exists", addedUser.Id)))
			return nil
		}
	}

	users = append(users, addedUser)
	jsonData, err := json.Marshal(&users)
	if err != nil {
		return err
	}

	_, err = f.Write(jsonData)
	if err != nil {
		return err
	}

	_, err = writer.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func performListCommand(args Arguments, writer io.Writer) error {
	f, ok := args[fileNameKey]
	if !ok {
		fmt.Errorf("-fileName flag has to be specified")
	}

	return printFile(f, writer)
}

func performRemoveCommand(args Arguments, writer io.Writer) error {
	fileName, ok := args[fileNameKey]
	if !ok || len(fileName) == 0 {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	id, ok := args[idKey]
	if !ok || len(id) == 0 {
		return fmt.Errorf("-id flag has to be specified")
	}

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var users []User
	if len(bytes) != 0 {
		err = json.Unmarshal(bytes, &users)
		if err != nil {
			return err
		}
	}

	founded := false
	foundedPos := -1
	for i, v := range users {
		if v.Id == id {
			founded = true
			foundedPos = i
			break
		}
	}

	if founded == false {
		writer.Write([]byte(fmt.Sprintf("Item with id %s not found", id)))
		return nil
	}

	if !(founded && foundedPos >= 0 && foundedPos < len(users)) {
		return fmt.Errorf("Unknown find error")
	}
	var newUsers []User
	newUsers = append(newUsers, users[:foundedPos]...)
	newUsers = append(newUsers, users[foundedPos+1:]...)

	jsonData, err := json.Marshal(&newUsers)
	if err != nil {
		return err
	}

	f.Close()
	os.Remove(fileName)

	newFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = newFile.Write(jsonData)
	if err != nil {
		return err
	}

	_, err = writer.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func performFindByIdCommand(args Arguments, writer io.Writer) error {
	fileName, ok := args[fileNameKey]
	if !ok || len(fileName) == 0 {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	id, ok := args[idKey]
	if !ok || len(id) == 0 {
		return fmt.Errorf("-id flag has to be specified")
	}

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var users []User
	if len(bytes) != 0 {
		err = json.Unmarshal(bytes, &users)
		if err != nil {
			return err
		}
	}

	founded := false
	foundedPos := -1
	for i, v := range users {
		if v.Id == id {
			founded = true
			foundedPos = i
			break
		}
	}

	if founded == false {
		writer.Write([]byte(""))
		return nil
	}

	if !(founded && foundedPos >= 0 && foundedPos < len(users)) {
		return fmt.Errorf("Unknown find error")
	}

	jsonData, err := json.Marshal(&users[foundedPos])
	if err != nil {
		return err
	}

	_, err = f.Write(jsonData)
	if err != nil {
		return err
	}

	_, err = writer.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func Perform(args Arguments, writer io.Writer) error {
	v, ok := args[operationKey]
	if !ok || len(v) == 0 {
		return fmt.Errorf("-operation flag has to be specified")
	}

	f, ok := args[fileNameKey]
	if !ok || len(f) == 0 {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	switch v {
	case "add":
		return performAddCommand(args, writer)
	case "list":
		return performListCommand(args, writer)
	case "remove":
		return performRemoveCommand(args, writer)
	case "findById":
		return performFindByIdCommand(args, writer)
	}

	return fmt.Errorf("Operation %s not allowed!", v)
}

func parseArgs() Arguments {
	operation := flag.String(operationKey, "", "operation type")
	item := flag.String(itemKey, "", "item")
	fileName := flag.String(fileNameKey, "", "file name for saving data")
	id := flag.String(idKey, "", "user id")

	flag.Parse()

	args := make(Arguments)
	args[operationKey] = *operation
	args[itemKey] = *item
	args[fileNameKey] = *fileName
	args[idKey] = *id

	fmt.Printf("Arguments %+v", args)

	return args
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
