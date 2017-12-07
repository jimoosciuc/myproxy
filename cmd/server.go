package main

import (
	"fmt"
	"os"
	"github.com/koko990/myproxy/cmd/app"
)

func main() {
	if err := app.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
