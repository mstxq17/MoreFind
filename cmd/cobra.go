package cmd

// Help template defines the format of the help message.
var helpTemplate = `{{.Long | trim}}

Usage:
  {{.CommandPath}} [flags]
{{if .Runnable}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

`

// Usage template defines the format of the usage message.
var usageTemplate = `Usage: {{.CommandPath}} [flags]`
