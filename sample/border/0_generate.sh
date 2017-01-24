#!/bin/bash

set -eux
cd `dirname $0`

: 'borderの値を変えてbitmap画像を複数生成します'
for border in `seq 0.0 0.1 2.0`; do
  target="b${border}.bmp"
  if [ -e "b${border}.png" ]; then
    continue
  fi
  ciedecompress -i "../src.jpg" -o $target -size 8 -border $border
done

: '生成されたbitマップをすべてjpegとpngに圧縮します'
for filename in $(ls *.bmp); do
  convert $filename ${filename//bmp/jpg}
  convert $filename ${filename//bmp/png}
done

: 'それぞれにPSNRとSSIMを計算します'
csvfile="9_quality.csv"
echo 'border,psnr,ssim' > $csvfile
for border in `seq 0.1 0.1 2.0`; do
  target="b${border}.png"
  psnr=$(imagediff -m psnr ../src.jpg $target)
  ssim=$(imagediff -m ssim ../src.jpg $target)
  echo "${border},${psnr},${ssim}" >> $csvfile
done

: '不要なbitmap画像を消去します'
rm *.bmp
