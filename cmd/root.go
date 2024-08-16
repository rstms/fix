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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/rstms/fix/vimfix"
)

var cfgFile string
var version = "1.0.5"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fix",
	Version: version,
	Short: "vim quickfix compiler output",
	Long: `
run a compile or lint command, scanning the output output for errors suitable for input to vim quickfix"
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	    //fleem = "foo"
	    os.Exit(vimfix.Fix(args[0], args[1:]...))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.Flags().BoolVarP(&vimfix.Quiet, "quiet", "q", false, "no echo stdout")
	rootCmd.Flags().BoolVarP(&vimfix.Verbose, "verbose", "v", false, "output diagnostics to stderr")
	rootCmd.Flags().BoolVarP(&vimfix.IgnoreStderr, "ignore-stderr", "E", false, "ignore stderr when scanning")
	rootCmd.Flags().BoolVarP(&vimfix.IgnoreStdout, "ignore-stdout", "O", false, "ignore stdout when scanning")
	rootCmd.Flags().BoolVarP(&vimfix.LocalizePaths, "localize", "l", false, "localize source file path in error output")
	rootCmd.Flags().BoolVarP(&vimfix.NoStripANSI, "no-strip", "S", false, "do not strip ANSI codes")
	rootCmd.Flags().BoolVarP(&vimfix.PrioritizeExitCode, "prioritize-exit-code", "x", false, "suppress prompt when command exits 0")
	rootCmd.Flags().StringVarP(&vimfix.ErrorFormat, "format", "f", "", "error format")
	rootCmd.Flags().StringVarP(&vimfix.OutputFile, "output", "o", "", "output to file")
	viper.BindPFlag("format", rootCmd.Flags().Lookup("format"))
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
