package version

import (
	"fmt"

	"github.com/fs714/github-star-manager/pkg/utils/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var StartCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func InitStartCmd() {
	StartCmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		command.Parent().PersistentFlags().VisitAll(func(flag *pflag.Flag) {
			flag.Hidden = true
		})
		command.Parent().HelpFunc()(command, strings)
	})
}

func printVersion() {
	versionInfo := `Version Information:
Version:		%s
Commit:			%s
Build Time:		%s
Go Version:		%s`

	fmt.Println(fmt.Sprintf(versionInfo, version.BaseVersion, version.GitVersion, version.BuildTime, version.GoVersion))
}
