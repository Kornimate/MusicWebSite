from pytube import YouTube
import os
import sys

DESTINATION = os.environ("DESTINATION_DIR")

def DownloadMusic(url):
    error_counter = 3
    while error_counter > 0:
        try:
            yt = YouTube(str(url))
            print(f"Downloading: {yt.title}")
            video = yt.streams.filter(only_audio=True).first()
            outfile = video.download(output_path=DESTINATION)
            base, ext = os.path.splitext(outfile)
            newfile = base + '.mp3'
            os.rename(outfile,newfile)
            print(f"{yt.title} downloaded")
            error_counter = 0
        except:
            if error_counter == 0:
                errorText = "file could not be downloaded"
                print(errorText)
                return errorText
            else:
                print(f"error while downloading file ({4-error_counter}. attempt)")
                error_counter -= 1
