FROM golang:1.23-alpine

WORKDIR /home

ENV HOST=:8080
ENV DOWNLOAD_PATH=/home/temp

COPY . .

RUN apk add py3-pip python3 ffmpeg

RUN pip install -r requirements.txt --break-system-packages

COPY ./python/cipher.py /usr/lib/python3.12/site-packages/pytube/cipher.py

RUN go mod download \ 
     && go mod tidy

RUN go build -o server .

EXPOSE 8080

CMD [ "./server" ]