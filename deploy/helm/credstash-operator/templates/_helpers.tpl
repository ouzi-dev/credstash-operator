{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "credstash-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "credstash-operator.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "credstash-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "credstash-operator.labels" -}}
helm.sh/chart: {{ include "credstash-operator.chart" . }}
{{ include "credstash-operator.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "credstash-operator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "credstash-operator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "credstash-operator.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "credstash-operator.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the rbac role to use
*/}}
{{- define "credstash-operator.rbacRoleName" -}}
{{- if .Values.rbac.create -}}
    {{ default (include "credstash-operator.fullname" .) .Values.rbac.roleName }}
{{- else -}}
    {{ default "default" .Values.rbac.roleName }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the rbac role binding to use
*/}}
{{- define "credstash-operator.rbacRoleBindingName" -}}
{{- if .Values.rbac.create -}}
    {{ default (include "credstash-operator.fullname" .) .Values.rbac.roleBindingName }}
{{- else -}}
    {{ default "default" .Values.rbac.roleBindingName }}
{{- end -}}
{{- end -}}

{{/*
Create the name of aws-credentials secret to use
*/}}
{{- define "credstash-operator.credentialsSecretName" -}}
{{- if .Values.awsCredentials.create -}}
    {{ default (printf "%s-%s" .Release.Name "aws-creds") .Values.awsCredentials.secretName }}
{{- else -}}
    {{ default "default" .Values.awsCredentials.secretName  }}
{{- end -}}
{{- end -}}