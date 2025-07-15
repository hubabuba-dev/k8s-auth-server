/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"kubectl-login/internal/cli"
)

func main() {
	err := cli.ExecuteCommand()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Login success.")
	}
}
