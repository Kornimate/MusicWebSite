from pytube import YouTube
import os
import sys

DESTINATION = '.'

error_counter = 3
while error_counter > 0:
    try:
        yt = YouTube(str(sys.argv[1]))
        print(f"Downloading: {yt.title}")
        video = yt.streams.filter(only_audio=True).first()
        outfile = video.download(output_path=DESTINATION)
        base, ext = os.path.splitext(outfile)
        newfile = base + '.mp3'
        os.rename(outfile,newfile)
        print(f"{yt.title} downloaded")
        error_counter = 0
    except IndexError:
        print("no input was given")
        error_counter = 0
    except:
        if error_count == 0:
            print("file could not be downloaded")
        else:
            print(f"error while downloading file ({4-error_counter}. attempt)")
            error_counter -= 1
