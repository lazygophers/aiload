package main

import ({{ if gt (len .RPCS) 0 }}
	"{{ $.PB.GoPackage }}"
	"{{ $.PB.GoPackage }}/internal/impl" {{ end }}
"{{ $.PB.GoPackage }}/internal/api"
)

var Routes = []*api.Route{ {{ range $key, $value := .RPCS }}
	{
	Method:  "{{ $value.Method }}",
	Path:    {{ $.PB.GoPackageName }}.RpcPath{{ $value.RpcName }},
	Handler: impl.ToHandler(impl.{{ $value.RpcName }}, "{{ with $value.Role }}{{ $value.Role }}{{ else }}user{{ end }}"),
	Role: "{{ with $value.Role }}{{ $value.Role }}{{ else }}public{{ end }}",
	},{{ end }}
}
