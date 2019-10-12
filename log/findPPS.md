run `ffmpeg -i video.mov -vcodec copy -vbsf h264_mp4toannexb -an  outfile.h264`
look for sps and pps in raw nalus, compare with fdsc extension and boom.

--> raw pps nalu from h264 `27640033 AC568047 0133E69E 6E020202 04`
--> fdsc extn `01640033 FFE10011 (27640033 AC568047 0133E69E 6E020202 04)010004 (28EE3CB0 FDF8F800)` 1. pps, 2. sps 