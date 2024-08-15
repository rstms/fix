package vimfix

import (
	"bytes"
	"fmt"
	"github.com/rstms/fix/util"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var Quiet bool
var Verbose bool
var IgnoreStderr bool
var IgnoreStdout bool
var LocalizePaths bool
var NoStripANSI bool
var ErrorFormat string
var OutputFile string
var PrioritizeExitCode bool

func stripANSI(s string) string {
	r, err := regexp.Compile("\x1B[@-_][0-?]*[ -/]*[@-~]")
	cobra.CheckErr(err)
	ret := r.ReplaceAll([]byte(s), []byte{})
	return string(ret)
}

func stripCRLF(s string) string {
	return strings.TrimSpace(s)
}

func eFormat(error, detail string) string {
	split := strings.SplitAfterN(detail, "-->", 2)
	if len(split) > 1 {
		detail = split[1]
	}
	return fmt.Sprintf("%s%s", detail, error)
}

func osPath(path string) string {
	dir, file := filepath.Split(path)
	parts := strings.Split(dir, string(filepath.Separator))
	parts = append(parts, file)
	path = filepath.Join(parts...)
	return path
}

func fixPath(line string) string {
	parts := strings.SplitAfterN(line, ":", 2)
	path := parts[0]
	tail := parts[1]
	return fmt.Sprintf("%s:%s", osPath(path), tail)
}

func forgeErrors(lines []string) []string {
	elines := []string{}
	failed := false

	eline := ""
	for _, line := range lines {
		words := strings.Fields(line)
		if failed {
			if eline != "" {
				eline := eFormat(eline, line)
				elines = append(elines, eline)
				eline = ""
			} else {
				if len(words) > 0 {
					if words[0] == "Error" || words[0] == "Warning" {
						eline = line
					}
				}
			}
		}
		if strings.Index(line, "Compiler run failed") != -1 {
			failed = true
		}
	}
	return elines
}

func isError(line string) bool {
	if line == "" {
		return false
	}
	if line[0] == '#' {
		return false
	}
	if strings.Contains(line, ":") {
		return true
	}
	return false
}

func genericErrors(lines []string) []string {
	elines := []string{}

	for _, line := range lines {
		if isError(line) {
			elines = append(elines, line)
		}
	}
	return elines
}

func blackErrors(lines []string) []string {
	re, err := regexp.Compile("^error: cannot format\\s(.*)")
	cobra.CheckErr(err)
	if Verbose {
		fmt.Fprintf(os.Stderr, "BEGIN blackErrors\n")
	}

	elines := []string{}

	for _, line := range lines {

		if Verbose {
			fmt.Fprintf(os.Stderr, "line: '%s'\n", line)
		}

		subs := re.FindStringSubmatch(line)
		if subs != nil && len(subs) > 1 {
			if Verbose {
				fmt.Fprintf(os.Stderr, "  subs: %v\n", subs)
			}
			line = subs[1]
			parts := strings.Split(line, ":")
			if Verbose {
				fmt.Fprintf(os.Stderr, "  parts: %v\n", parts)
			}
			file := strings.TrimSpace(parts[0])
			emsg := strings.TrimSpace(parts[1])
			row, err := strconv.Atoi(strings.TrimSpace(parts[2]))
			cobra.CheckErr(err)
			col, err := strconv.Atoi(strings.TrimSpace(parts[3]))
			cobra.CheckErr(err)
			source := strings.TrimSpace(parts[4])

			line = fmt.Sprintf("%s:%d:%d: [black] %s %s", file, row, col, emsg, source)
			if Verbose {
				fmt.Fprintf(os.Stderr, "    file: '%s'\n", file)
				fmt.Fprintf(os.Stderr, "    emsg: '%s'\n", emsg)
				fmt.Fprintf(os.Stderr, "     row: '%d'\n", row)
				fmt.Fprintf(os.Stderr, "     col: '%d'\n", col)
				fmt.Fprintf(os.Stderr, "  source: '%s'\n", source)
				fmt.Fprintf(os.Stderr, "   eline: '%s'\n", line)
			}
			elines = append(elines, line)
		}
	}
	if Verbose {
		fmt.Fprintf(os.Stderr, "END blackErrors\n")
	}
	return elines
}

// rgrep -Hn installation_timeout
//
// roles/vmware-instance/tasks/netboot_wait.yml:18:  async: "{{ installation_timeout }}"

func rgrepErrors(lines []string) []string {
	elines := []string{}
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) > 2 {
			file := parts[0]
			row := parts[1]
			col := 1
			pos := len(file) + len(row) + 2
			msg := line[pos:]
			line = fmt.Sprintf("%s:%s:%d:%s", file, row, col, msg)
			elines = append(elines, line)
		}
	}
	return elines
}

func tryQuickfix(elines []string) int {
	confirmed, err := util.Confirm("fix")
	cobra.CheckErr(err)
	if confirmed {
		quickfixFile := ".quickfix"
		if LocalizePaths {
			flines := []string{}
			for _, eline := range elines {
				flines = append(flines, fixPath(eline))
			}
			elines = flines
		}
		err := os.WriteFile(quickfixFile, []byte(strings.Join(elines, "\n")), 0600)
		cobra.CheckErr(err)

		editor, ok := os.LookupEnv("EDITOR")
		if !ok {
			editor, ok = os.LookupEnv("VISUAL")
			if !ok {
				editor = "vim"
			}
		}
		fmt.Printf("calling %s...\n", editor)

		cmd := exec.Command(editor, "-n", "-q", quickfixFile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		cobra.CheckErr(err)
		err = os.Remove(quickfixFile)
		cobra.CheckErr(err)
	}
	return -1
}

var formats = map[string]func([]string) []string{
	"forge":   forgeErrors,
	"black":   blackErrors,
	"rgrep":   rgrepErrors,
	"generic": genericErrors,
}

func getFormatter(command string) func([]string) []string {

	formatter := "generic"

	// override command with --format option
	if ErrorFormat != "" {
		command = ErrorFormat
	}

	fmtFunc, ok := formats[command]
	if ok {
		formatter = command
	} else {
		formatter = "generic"
		fmtFunc = formats[formatter]
	}

	if Verbose {
		fmt.Fprintf(os.Stderr, "command: %s, formatter: %s\n", command, formatter)
	}
	return fmtFunc
}

func strippedLines(buf []byte, stripper func(string) string) []string {
	lines := strings.Split(string(buf), "\n")
	stripped := []string{}
	for _, line := range lines {
		stripped = append(stripped, stripper(line))
	}
	return stripped
}

func customizeArgs(command string, args ...string) []string {
	switch command {
	case "rgrep":
		if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
			args = append([]string{"-Hn"}, args...)
		}
	}
	return args
}

func Fix(command string, args ...string) int {

	args = customizeArgs(command, args...)

	cmd := exec.Command(command, args...)
	var obuf bytes.Buffer
	var ebuf bytes.Buffer
	cmd.Stdout = &obuf
	cmd.Stderr = &ebuf

	err := cmd.Run()
	var exitCode int
	exiterr, ok := err.(*exec.ExitError)
	if ok {
		exitCode = exiterr.ExitCode()
	} else {
		cobra.CheckErr(err)
	}

	if !Quiet {
		_, err := os.Stdout.Write(obuf.Bytes())
		cobra.CheckErr(err)
	}

	_, err = os.Stderr.Write(ebuf.Bytes())
	cobra.CheckErr(err)

	var stripper func(string) string
	if NoStripANSI {
		stripper = stripCRLF
	} else {
		stripper = stripANSI
	}

	formatter := getFormatter(command)

	elines := []string{}
	if !IgnoreStdout {
		stripped := strippedLines(obuf.Bytes(), stripper)
		formatted := formatter(stripped)
		elines = append(elines, formatted...)
	}

	if !IgnoreStderr {
		stripped := strippedLines(ebuf.Bytes(), stripper)
		formatted := formatter(stripped)
		elines = append(elines, formatted...)
	}

	for {
		if PrioritizeExitCode && exitCode == 0 {
			//fmt.Printf("exitCode=%v\n", exitCode)
			break
		}

		if len(elines) > 0 {
			return tryQuickfix(elines)
		}
		break
	}

	if OutputFile != "" && exitCode == 0 {
		// write output file only if no errors
		err := os.WriteFile(OutputFile, obuf.Bytes(), 0600)
		cobra.CheckErr(err)
	}

	return exitCode
}
