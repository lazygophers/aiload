package {{ .PB.GoPackageName }}

const ({{ range $key, $value := .RPCS }}
RpcPath{{ $value.RpcName }} = "/{{ TrimPrefix $value.Path "/" }}"{{ end }}
)
