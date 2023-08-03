{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "alidns-webhook.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "alidns-webhook.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "alidns-webhook.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "alidns-webhook.labels" -}}
helm.sh/chart: {{ include "alidns-webhook.chart" . }}
{{ include "alidns-webhook.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "alidns-webhook.selectorLabels" -}}
app.kubernetes.io/name: {{ include "alidns-webhook.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "alidns-webhook.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "alidns-webhook.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}


{{- define "alidns-webhook.selfSignedIssuer" -}}
{{ printf "%s-selfsign" (include "alidns-webhook.fullname" .) }}
{{- end -}}

{{- define "alidns-webhook.rootCAIssuer" -}}
{{ printf "%s-ca" (include "alidns-webhook.fullname" .) }}
{{- end -}}

{{- define "alidns-webhook.rootCACertificate" -}}
{{ printf "%s-ca" (include "alidns-webhook.fullname" .) }}
{{- end -}}

{{- define "alidns-webhook.servingCertificate" -}}
{{ printf "%s-tls" (include "alidns-webhook.fullname" .) }}
{{- end -}}
