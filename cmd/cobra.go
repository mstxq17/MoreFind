package cmd

// Help template defines the format of the help message.
var helpTemplate = `{{.Short| trim}}
{{.Long| trim}}

Usage:
  {{.CommandPath}} [params] [flags]
{{if .Runnable}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

Global Flags:
  -o, --output string                                            Specifies the output file path.
`

// Usage template defines the format of the usage message.
var usageTemplate = `Usage: {{.CommandPath}} [flags]`

var deduHelpTemplate = `{{.Short| trim}}
{{.Long| trim}}

Usage:
  {{.CommandPath}} [flags]
{{if .Runnable}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

Global Flags:
  -o, --output string                                            Specifies the output file path.
`
