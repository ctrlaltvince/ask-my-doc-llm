{{- define "backend.name" -}}
backend
{{- end }}

{{- define "backend.fullname" -}}
{{ .Release.Name }}-backend
{{- end }}
