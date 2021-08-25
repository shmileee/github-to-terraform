terraform import github_repository_collaborator.{{ .Repository.Name}}-{{ .Collaborator.Login }} {{ .Repository.Name}}:{{ .Collaborator.Login }}
