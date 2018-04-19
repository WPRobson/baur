package cmd

import (
	"os"
	"path"
	"strings"

	"github.com/simplesurance/sisubuild/cfg"
	"github.com/simplesurance/sisubuild/sblog"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(appInitCmd)
}

const appInitLongHelp = `
Create an application config file in the current directory.
The name parameter is set to the current directory name.`

var appInitCmd = &cobra.Command{
	Use:   "appinit",
	Short: "creates an application config file in the current directory",
	Long:  strings.TrimSpace(appInitLongHelp),
	Run:   appInit,
}

func appInit(cmd *cobra.Command, args []string) {
	mustFindRepositoryRoot()

	cwd, err := os.Getwd()
	if err != nil {
		sblog.Fatal(err)
	}
	appName := path.Base(cwd)

	err = cfg.NewApplicationFile(appName, path.Join(cwd, cfg.ApplicationFile))
	if err != nil {
		if os.IsExist(err) {
			sblog.Fatalf("%s already exist", cfg.ApplicationFile)
		}

		sblog.Fatal(err)
	}

	sblog.Infof("written application configuration file to %s",
		cfg.ApplicationFile)
}
