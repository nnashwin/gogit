package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

func main() {
	cli.NewApp().Run(os.Args)
	app := cli.NewApp()
	app.Name = "gogit"
	app.Usage = "Letting your manage all of our Github accounts in style"
	app.Action = func(c *cli.Context) error {
		fmt.Println("Go Get Github")
		return nil
	}

	app.Run(os.Args)
}
