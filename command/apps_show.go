package command

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/simplesurance/baur/format"
	"github.com/simplesurance/baur/format/csv"
	"github.com/simplesurance/baur/format/table"
	"github.com/simplesurance/baur/log"
)

type appsShowConf struct {
	csv bool
}

var appsShowCmd = &cobra.Command{
	Use:   "show <APP-NAME>|<PATH>",
	Short: "show configuration of an application",
	Run:   appsShow,
	Args:  cobra.ExactArgs(1),
}

var appsShowConfig appsShowConf

func init() {
	appsShowCmd.Flags().BoolVar(&appsShowConfig.csv, "csv", false,
		"Show output in RFC4180 CSV format")

	appsCmd.AddCommand(appsShowCmd)
}

func appsShow(cmd *cobra.Command, args []string) {
	var formatter format.Formatter

	repo := MustFindRepository()
	app := mustArgToApp(repo, args[0])

	if appsShowConfig.csv {
		formatter = csv.New(nil, os.Stdout)
	} else {
		formatter = table.New(nil, os.Stdout)
	}

	mustWriteRow(formatter, []interface{}{fmtVertTitle(appsShowConfig.csv, "Name"), app.Name})
	mustWriteRow(formatter, []interface{}{fmtVertTitle(appsShowConfig.csv, "Path"), app.RelPath})
	mustWriteRow(formatter, []interface{}{fmtVertTitle(appsShowConfig.csv, "Build Command"), app.BuildCmd})

	if len(app.Outputs) != 0 {
		mustWriteRow(formatter, []interface{}{})
		mustWriteRow(formatter, []interface{}{fmtVertTitle(appsShowConfig.csv, "Outputs")})

		for i, art := range app.Outputs {
			mustWriteRow(formatter, []interface{}{fmtVertTitle(appsShowConfig.csv, "Local"), art.String()})
			mustWriteRow(formatter, []interface{}{fmtVertTitle(appsShowConfig.csv, "Remote"), art.UploadDestination()})

			if i+1 < len(app.Outputs) {
				mustWriteRow(formatter, []interface{}{})
			}
		}
	}

	if len(app.BuildInputPaths) != 0 {
		mustWriteRow(formatter, []interface{}{})
		mustWriteRow(formatter, []interface{}{fmtVertTitle(appsShowConfig.csv, "Inputs")})

		for _, bi := range app.BuildInputPaths {
			mustWriteRow(formatter, []interface{}{fmtVertTitle(appsShowConfig.csv, bi.Type()), bi.String()})
		}

	}

	if err := formatter.Flush(); err != nil {
		log.Fatalln(err)
	}
}
