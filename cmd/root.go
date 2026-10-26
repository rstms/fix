/*
Copyright © 2024 Matt Krueger <mkrueger@rstms.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/rstms/fix/vimfix"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var version = "1.0.10"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "fix [[OPTS] --] COMMAND [COMMAND_OPTS] [COMMAND_ARGS...]",
	Version: version,
	Short:   "vim quickfix exec",
	Long: `
execute a command and offer to run vim in quickfix mode against any generated errors
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vimfix.Quiet = ViperGetBool("quiet")
		vimfix.Verbose = ViperGetBool("verbose")
		vimfix.IgnoreStderr = ViperGetBool("ignore_stderr")
		vimfix.IgnoreStdout = ViperGetBool("ignore_stdout")
		vimfix.LocalizePaths = ViperGetBool("localize")
		vimfix.NoStripANSI = ViperGetBool("no_strip")
		vimfix.PrioritizeExitCode = ViperGetBool("prioritize_exit_code")
		vimfix.ErrorFormat = ViperGetString("error_format")
		vimfix.OutputFile = ViperGetString("output_file")
		os.Exit(vimfix.Fix(args[0], args[1:]...))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	CobraInit(rootCmd)
	OptionSwitch(rootCmd, "quiet", "q", "no echo stdout")
	OptionSwitch(rootCmd, "ignore-stderr", "E", "ignore stderr when scanning")
	OptionSwitch(rootCmd, "ignore-stdout", "O", "ignore stdout when scanning")
	OptionSwitch(rootCmd, "localize", "l", "localize source file path in error output")
	OptionSwitch(rootCmd, "no-strip", "S", "do not strip ANSI codes")
	OptionSwitch(rootCmd, "prioritize-exit-code", "x", "suppress prompt when command exits 0")
	OptionString(rootCmd, "format", "f", "", "error format")
	OptionString(rootCmd, "output", "o", "", "output to file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".fix" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".fix")
	}
	viper.SetEnvPrefix("fix")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
