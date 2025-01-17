{{- define "vectorEnv" }}
- name: VECTOR_SELF_NODE_NAME
  valueFrom:
    fieldRef:
      fieldPath: spec.nodeName
- name: VECTOR_SELF_POD_NAME
  valueFrom:
    fieldRef:
      fieldPath: metadata.name
- name: VECTOR_SELF_POD_NAMESPACE
  valueFrom:
    fieldRef:
      fieldPath: metadata.namespace
{{- end }}

{{- define "vectorMounts" }}
- name: vector-data-dir
  mountPath: "/vector-data"
- name: vector-config-dir
  mountPath: /etc/vector/dynamic
- name: vector-sample-config-dir
  mountPath: /etc/vector/default/defaults.json
  readOnly: true
  subPath: defaults.json
{{- end }}
