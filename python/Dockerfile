FROM python:3.12-alpine

ENV DESTINATION_DIR=/home
ENV HOST=0.0.0.0

WORKDIR /app

COPY . .

RUN pip3 install -R requirements.txt

EXPOSE 8080

CMD ["python3","server.py"]