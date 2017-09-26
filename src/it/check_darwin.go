package main

import "os"
import "os/exec"
import "fmt"

var git = "git"
var goc = "go"
var ext = ""

func AptGetInstall(name string) error {
	err := exec.Command("which", "apt-get").Run()
	if err != nil {
		return err
	}

	fmt.Println("Please input your password to install ", name)
	install := exec.Command("sudo", "apt-get", "install", "golang")
	install.Stdin = os.Stdin
	install.Stdout = os.Stdout
	err = install.Run()
	return err
}

func SystemChecks() {
	var err error

	err = exec.Command("which", "go").Run()
	if err != nil {
		fmt.Println("The Go programming language is required for IT and will now be installed...")
		err := AptGetInstall("golang")
		if err != nil {
			fmt.Println("The Go programming language cannot be automatically installed on your system.")
			fmt.Println("Please visit http://golang.org/dl/ and install it manually.")
			os.Exit(1)
		}
	}
	
}
