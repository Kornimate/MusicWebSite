from flask import Flask, request, Response, send_file
from flask_cors import CORS
import os
import python.service as service

app = Flask(__name__)
cors = CORS(app)


@app.route('/api/v1/music', methods=['POST'])
def Music():
    if request.method != "POST":
        return Response("{'success':'false', 'content':'invalid method'}",400,mimetype='application/json')
    
    url = request.form['url']
    file_data = service.Download(url)
    
    if file_data.success:
        return send_file(file_data.path,mimetype="application/octet-stream", download_name=file_data.download_name)


if __name__ == "__main__":
    
    HOST = os.environ("HOST")
    
    if HOST == "" or HOST == None:
        HOST = "localhost"
    
    app.run(host="0.0.0.0",port=8080)