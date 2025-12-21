FROM golang:1.25-alpine

RUN go install github.com/air-verse/air@latest && \
  go install github.com/go-delve/delve/cmd/dlv@latest \
  go install github.com/a-h/templ/cmd/templ@latest

RUN apk add --no-cache git gcc musl-dev

RUN echo fs.inotify.max_user_watches=524288 | tee -a /etc/sysctl.conf && \
    echo fs.inotify.max_user_instances=512 | tee -a /etc/sysctl.conf
    
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENTRYPOINT ["air","-c",".air.debug.toml"]
