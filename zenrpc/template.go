package main

import (
	"text/template"

	"github.com/IlyushaZ/zenrpc/v3/parser"
)

var (
	serviceTemplate = template.Must(template.New("service").
		Funcs(template.FuncMap{"definitions": parser.Definitions}).
		Parse(`
// Code generated by modified version of zenrpc; DO NOT EDIT.
package {{.PackageName}}
import (
	"encoding/json"
	"context"

	"github.com/IlyushaZ/zenrpc/v3"
	"github.com/IlyushaZ/zenrpc/v3/smd"
	{{ range .ImportsForGeneration}}
		{{if .Name}}{{.Name.Name}} {{end}}{{.Path.Value}}
	{{- end }}
)
var RPC = struct {
{{ range .Services}}
	{{.Name}} struct { {{range $i, $e := .Methods }}{{if $i}}, {{end}}{{.Name}}{{ end }} string } 
{{- end }}
}{	
	{{- range .Services}}
		{{.Name}}: struct { {{range $i, $e := .Methods }} {{if $i}}, {{end}}{{.Name}}{{ end }} string }{ 
			{{- range .Methods }}
				{{.Name}}:   "{{.LowerCaseName}}",
			{{- end }}
		}, 	
	{{- end }}
}
{{ range $s := .Services}}
	func ({{.Name}}) SMD() smd.ServiceInfo {
		return smd.ServiceInfo{}
	}

	// Invoke is as generated code from zenrpc cmd
	func (s {{.Name}}) Invoke(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
		resp := zenrpc.Response{}
		switch method { 
		{{- range .Methods }}
			case RPC.{{$s.Name}}.{{.Name}}: {{ if .Args }}
					{{- $arg := index .Args 0 }}
					var {{ $arg.Name }} {{ $arg.Type }}

					if len(params) > 0 {
						if err := json.Unmarshal(params, &{{ $arg.Name }}); err != nil {
							return zenrpc.NewResponseError(nil, zenrpc.InvalidParams, "", err.Error())
						}
					}

					{{ range .DefaultValues }}
						{{.Comment}}
						if args.{{.CapitalName}} == nil {
							var v {{.Type}} = {{.Value}}
							args.{{.CapitalName}} = &v
						}
					{{ end }}
				{{ end }} {{if .Returns}}
					resp.Set(s.{{.Name}}({{if .HasContext}}ctx, {{end}} {{ range .Args }}{{if and (not .HasStar) .HasDefaultValue}}*{{end}}{{.Name}}, {{ end }}))
				{{else}}
					s.{{.Name}}({{if .HasContext}}ctx, {{end}} {{ range .Args }}{{if and (not .HasStar) .HasDefaultValue}}*{{end}}args.{{.CapitalName}}, {{ end }})
				{{end}}
		{{- end }}
		default:
			resp = zenrpc.NewResponseError(nil, zenrpc.MethodNotFound, "", nil)
		}

		return resp
	}
{{- end }}
`))
)
