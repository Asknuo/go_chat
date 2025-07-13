package flag

import (
	"fmt"
	"gochat/global"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	sqlFlag = &cli.BoolFlag{
		Name:  "sql",
		Usage: "Initialize SQL database",
	}
)

func Run(c *cli.Context) error {
	if c.NumFlags() > 1 {
		err := cli.NewExitError("Too many flags provided", 1)
		if err != nil {
			global.Log.Error("Error in Run function")
			return err
		}
	}
	switch {
	case c.Bool(sqlFlag.Name):
		err := SqlMigrate()
		if err != nil {
			global.Log.Error("表结构迁移失败:")
			return err
		} else {
			global.Log.Info("表结构迁移成功")
		}
	}
	return nil
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "gochat"
	app.Flags = []cli.Flag{
		sqlFlag,
	}
	app.Action = Run
	return app
}

func InitFlag() {
	if len(os.Args) > 1 {
		app := NewApp()
		err := app.Run(os.Args)
		if err != nil {
			global.Log.Error("Failed to run app")
			os.Exit(1)
		} else {
			global.Log.Info("App ran successfully")
		}
		if os.Args[1] == "-h" || os.Args[1] == "-help" {
			fmt.Println("Displaying help message...")
		}
		os.Exit(0)
	}
}
