package mysqlmd

import (
	"strings"

	"github.com/YeHeng/go-web-api/pkg/cmd/mysqlmd"

	"github.com/spf13/cobra"
)

var (
	dbAddr    string
	dbUser    string
	dbPass    string
	dbName    string
	genTables string
)

func NewCmdMySQLMd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mysqlmd -a [addr] -u [username] -p [password] -n [schema name] -t [table names]",
		Short: "Generate database Markdown Document",
		Annotations: map[string]string{
			"IsActions": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {

			dbName = strings.ToLower(dbName)
			genTables = strings.ToLower(genTables)

			mysqlmd.RunMySQLMD(dbAddr, dbUser, dbPass, dbName, genTables)

		},
	}

	cmd.Flags().StringVarP(&dbAddr, "addr", "a", "127.0.0.1:3306", "mysql address，example：127.0.0.1:3306")
	cmd.Flags().StringVarP(&dbUser, "user", "u", "root", "mysql username")
	cmd.Flags().StringVarP(&dbPass, "pass", "p", "", "mysql password")
	cmd.Flags().StringVarP(&dbName, "name", "n", "", "database name")
	cmd.Flags().StringVarP(&genTables, "table", "t", "*", "tables names, default: *, split with \",\"")

	return cmd
}
