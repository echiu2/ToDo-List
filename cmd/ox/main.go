package main

import (
	"context"
	"fmt"
	"os"

	"todolist"
	_ "todolist/app/tasks"
	_ "todolist/app/models"

	"github.com/wawandco/oxpecker/cli"
	"github.com/wawandco/oxpecker/tools/soda"
)

// main function for the tooling cli, will be invoked by Oxpecker
// when found in the source code. In here you can add/remove plugins that
// your app will use as part of its lifecycle.
func main() {
	cli.Use(soda.Plugins(todolist.Migrations)...)
	err := cli.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Printf("[error] %v \n", err.Error())

		os.Exit(1)
	}
}