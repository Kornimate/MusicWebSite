from pytube import YouTube
import os
import sys

def DownloadMusic(url, dir):
    error_counter = 3
    while error_counter > 0:
        try:
            yt = YouTube(str(url))
            print(f"Downloading: {yt.title}")
            video = yt.streams.filter(only_audio=True).first()
            outfile = video.download(output_path=dir)
            base, ext = os.path.splitext(outfile)
            newfile = base + '.mp3'
            os.rename(outfile,newfile)
            print(f"{yt.title} downloaded to {dir}")
            error_counter = 0
            if not os.path.exists(newfile):
                print(f"{newfile} does not exists")
                exit(1)
            else:
                print(f"{newfile} exists")
        except Exception as e:
            print(str(e))
            print(repr(e))
            if error_counter == 0:
                errorText = "file could not be downloaded"
                print(errorText)
                exit(1)
            else:
                print(f"error while downloading file ({4-error_counter}. attempt)")
                error_counter -= 1

if __name__ == "__main__":
    DownloadMusic(sys.argv[1],sys.argv[2])
