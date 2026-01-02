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
{{- end -}}
