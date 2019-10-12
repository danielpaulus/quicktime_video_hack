run `ffmpeg -i video.mov -vcodec copy -vbsf h264_mp4toannexb -an  outfile.h264`
look for sps and pps in raw nalus, compare with fdsc extension and boom. 