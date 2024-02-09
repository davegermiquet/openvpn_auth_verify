package main

import (
		"database/sql"
		"fmt"
		"golang.org/x/crypto/bcrypt"
		"github.com/akamensky/argparse"
		"os"
		"errors"
		_ "github.com/mattn/go-sqlite3"
	)

  const conf_location string = "user.db"

  type payload_struct struct {
	username string
	password string
	errorMsg string
	exitLevel int
}

  type command_struct struct {
	command func(db *sql.DB,payload payload_struct) int
	payload payload_struct
  }

  func FileExists(filePath string) (bool, error) {
    info, err := os.Stat(filePath)
    if err == nil {
        return !info.IsDir(), nil
    }
    if errors.Is(err, os.ErrNotExist) {
        return false, nil
    }
    return false, err
}

  func does_database_exist(payload payload_struct) int  {
	exist, err := FileExists(conf_location)
	if !(exist || err != nil ) {
		payload.errorMsg = "Database file not found or not created"
		payload.exitLevel = 5
		return exit_with_error(nil,payload)
	}
	return 0
  }

  func user_exists(db *sql.DB, username string) bool{

    var count int

	statement,_:= db.Prepare("SELECT count(*) FROM user where user = ?")
	
    row , _ := statement.Query(username)
	row.Next()
	row.Scan(&count)
	row.Close()
	if  count == 1 {
		return true
	} else {
		return false
	}
  }

  var exit_with_error = func( db *sql.DB, payload payload_struct) int {
	 fmt.Println(payload.errorMsg)
	 return payload.exitLevel
  }

  var del_user = func( db *sql.DB, payload payload_struct) int {

	exists:= does_database_exist(payload) 
	if exists != 0 {
		return exists
	}
	

	if  (user_exists(db,payload.username)) {
		query:="delete FROM user WHERE user = ?" 
	    stat,_:= db.Prepare(query)
		stat.Exec(payload.username)
	} else {
		fmt.Println("User doesn't exist")
	}
	return 1
  }

  var list_user = func( db *sql.DB, payload payload_struct) int {

	exists:= does_database_exist(payload) 
	if exists != 0 {
		return exists
	}

	statement,_:= db.Prepare("SELECT user FROM user")

	rows, _ := statement.Query()
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			break
		}
		fmt.Println(name)
	}
	
	return 1
  }

  var create_database = func( db *sql.DB, payload payload_struct) int{
	statement, _:= db.Prepare("CREATE TABLE IF NOT EXISTS USER(user VARCHAR UNIQUE, password BLOB)")
	statement.Exec()
	return 1
}

var add_user = func( db *sql.DB, payload payload_struct) int {
	
	exists:= does_database_exist(payload) 
	if exists != 0 {
		return exists
	}

	if !(user_exists(db,payload.username)) {
		hash_password, _:= bcrypt.GenerateFromPassword([]byte(payload.password),2)
		statement, _:= db.Prepare("INSERT INTO USER(user,password) VALUES(?,?)")
		statement.Exec(payload.username,hash_password)
	} else {
		fmt.Println("User already exist")
	}

	return 1
}


var logon_user  = func (db *sql.DB,payload  payload_struct) int {

	exists:= does_database_exist(payload) 
	if exists != 0 {
		return exists
	}

	statement,_:= db.Prepare("SELECT password FROM user where user = ?")
	var hash_password []byte
	username:= os.Getenv("username")
	password:= os.Getenv("password")
    row , err:= statement.Query(username)
	if err != nil {
		payload.exitLevel = 8
		payload.errorMsg = "Database issue: " + err.Error()
		return exit_with_error (nil,payload)
	}
	row.Next()
	row.Scan(&hash_password)
	row.Close()
	error:=bcrypt.CompareHashAndPassword(hash_password,[]byte(password))
	if error == nil {
		fmt.Println("Success")
		return 0
	}
	return 1
}


func parse_arguments() command_struct {

	var command command_struct
 
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	// Create string flag
	username := parser.String("u", "user", &argparse.Options{Required: false, Help: "Username for action specified"})
	password := parser.String("p", "password", &argparse.Options{Required: false, Help: "Password for action specified"})
	createDatabase:= parser.Flag("c", "create", &argparse.Options{Required: false, Help: "Create Database"})
	addUser:= parser.Flag("a", "add", &argparse.Options{Required: false, Help: "Enable Add User Action"})
	delUser:= parser.Flag("d", "delete", &argparse.Options{Required: false, Help: "Enable Delete User Action"})
	listUser:= parser.Flag("l", "list", &argparse.Options{Required: false, Help: "List Users"})

	// Parse input
	err := parser.Parse(os.Args)
	if *createDatabase {
		command.command = create_database
	} else if *addUser {
		command.command = add_user
		if len(*username) > 3 &&  len(*password) > 5 {
			command.payload.username = *username
			command.payload.password = *password
		} else {
			command.command = exit_with_error
			command.payload.errorMsg = "Sorry you need username nad password for this option or password needs to be more than 6 or more characters"
			command.payload.exitLevel = 5
		}
	} else if *delUser {
		command.command = del_user
		if len(*username) > 3 {
			command.payload.username = *username
		} else {
			command.payload.errorMsg = "Sorry you need username for this option"
			command.command  = exit_with_error
			command.payload.exitLevel = 5
		}
	} else if *listUser{
		command.command = list_user
		} else {
		command.command = logon_user
	}

	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		command.payload.exitLevel = 5
		command.payload.errorMsg = "Invalid Usage"
		command.command = exit_with_error
	}
	// Finally print the collected string
	return command
}

func main() {

	// change this to location of where you want the user.db to be store on the openvpn server (/etc/openvpn/server/user.db)
	var result int
	db, _ :=  sql.Open("sqlite3", conf_location)
	var command command_struct = parse_arguments()
	if !(command.command == nil){
		result = command.command(db,command.payload)
	}
	db.Close()
	os.Exit(result)
}