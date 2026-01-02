{{/*
Expand the name of the chart.
*/}}
{{- define "app.name" -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "app.fullname" -}}
{{- printf "%s-%s" ( include "app.name" . ) .Values.nameSuffix | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{/*
* TODO: Write comment.
*/}}
{{- define "app.global.env" -}}
{{- if not .Values.global.env }}
{{- fail "'global.env' parameter is required but not defined!" }}
{{- end }}
{{- if not (contains .Values.global.env .Release.Namespace) }}
{{- fail "the namespace must contain 'global.env'" }}
{{- end }}
{{- .Values.global.env }}
{{- end -}}

{{/*
Common labels.
*/}}
{{- define "app.labels" -}}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{ include "app.selectorLabels" . }}
{{- end -}}

{{/*
Selector labels.
*/}}
{{- define "app.selectorLabels" -}}
app.kubernetes.io/name: {{ include "app.fullname" . }}
app.kubernetes.io/env: {{ include "app.global.env" . }}
{{- end -}}
