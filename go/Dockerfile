FROM golang:1.23-alpine

WORKDIR /app

ENV HOST=:8080

RUN apk add ffmpeg py3-pip

RUN pip install yt-dlp --break-system-package

COPY . .

RUN go mod download && go mod tidy

RUN go build -o server .

EXPOSE 8080

CMD ["./server"]