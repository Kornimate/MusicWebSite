FROM golang:1.23-alpine

WORKDIR /home

ENV HOST=:8080
ENV DOWNLOAD_PATH=/tmp

RUN apk add ffmpeg py3-pip

RUN pip install -U yt-dlp --break-system-package

COPY . .

COPY sshd_config /etc/ssh/

RUN go mod download \ 
     && go mod tidy

RUN go build -o server .
     
RUN apk add openssh \
    && echo "root:Docker!" | chpasswd \
    && chmod +x /home/start_container.sh \
    && cd /etc/ssh/ \
    && ssh-keygen -A \
    && cd /home

EXPOSE 8080 2222

ENTRYPOINT [ "/home/start_container.sh" ]