package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/howeyc/gopass"
)

func PromptUnix(user User) User {
    fmt.Printf("Username: ")
    fmt.Scan(&user.Username)
    fmt.Printf("Password: ")
    password, _ := gopass.GetPasswd()
    user.Password = string(password)
    return user
}

func PromptWindows(user User) User {
    fmt.Printf("Username: ")
    fmt.Scan(&user.Username)
    fmt.Printf("Password: ")
    fmt.Scan(&user.Password)
    return user
}

func Prompt(user User) User {
	var err error

	switch runtime.GOOS {
	case "linux":
        PromptUnix(user)
	case "windows":
        PromptWindows(user)
	case "darwin":
        PromptUnix(user)
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
    return user
}

func CheckCredentials(user User, path string, fpath string) bool {
    // if file doesn't exist then create one
    if _, err := os.Stat(fpath); os.IsNotExist(err) {
        fmt.Println("credential.json doesn't exist! creating file...")

        os.MkdirAll(path, os.ModePerm)

        user = Prompt(user)

    // if file exist then check credential
    } else {
        file, _ := ioutil.ReadFile(fpath)
        if err := json.Unmarshal([]byte(file), &user); err != nil {
            log.Fatal(err)
        }

        if user.Username == "" ||  user.Password == "" {
            user = Prompt(user)
        }
    }

    file, _ := json.MarshalIndent(user, "", " ")
    ioutil.WriteFile(fpath, file, 0644)

    return true
}

