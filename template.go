package main

var tmplDockerfile = `FROM {{.BaseImage}}:{{.BaseVersion}}
LABEL maintainer="edznux@gmail.com"

ARG NAME={{.Name}}

{{ if or (eq .BaseImage "alpine") (eq .BaseVersion "alpine") }}
RUN apk update && apk upgrade -y
RUN adduser -D ${NAME} && mkdir -p /home/${NAME}/workdir && chown -R ${NAME} /home/${NAME}
{{ else }}
RUN apt-get update && apt-get upgrade -y
RUN useradd -s /bin/bash -m ${NAME} && mkdir /home/${NAME}/workdir && chown ${NAME}:${NAME} /home/${NAME}/workdir
{{ end }}

USER ${NAME}
WORKDIR /home/${NAME}/workdir
ENTRYPOINT ["/home/{{.Name}}/"]
`

var tmplAlias = `#!/usr/bin/env bash

# Pass all arguments to the docker container
# We use the workdir subdirectory to avoid dealing with .bash_history and other crap
{{.Name}}(){
	docker run --rm -it \
		-v "$(pwd):/home/{{.Name}}/workdir" \
		$DOCKER_REPO_PREFIX/{{.Name}} \
		$@
}
`
