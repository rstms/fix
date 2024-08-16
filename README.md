# fix
```

execute a command and offer to run vim in quickfix mode against any generated errors

Usage:
  fix COMMAND [OPTS] [ARGS...] [flags]

Flags:
  -f, --format string          error format
  -h, --help                   help for fix
  -E, --ignore-stderr          ignore stderr when scanning
  -O, --ignore-stdout          ignore stdout when scanning
  -l, --localize               localize source file path in error output
  -S, --no-strip               do not strip ANSI codes
  -o, --output string          output to file
  -x, --prioritize-exit-code   suppress prompt when command exits 0
  -q, --quiet                  no echo stdout
  -v, --verbose                output diagnostics to stderr
      --version                version for fix
```
