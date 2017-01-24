#!/bin/bash

set -eux
cd `dirname $0`

: 'ciedecompressのかかった画像とかかっていない画像を用意します'
raw="../src.jpg"
com="../border/b1.0.png"

: '2枚の画像をそれぞれ5〜95のqualityでjpeg圧縮します'
for quality in `seq 5 5 90`; do
  convert $raw -quality $quality "raw_${quality}.jpg"
  convert $com -quality $quality "com_${quality}.jpg"
done

: 'それぞれにPSNRとSSIMを計算します'
csvfile="9_quality.csv"
echo 'quality,raw_psnr,com_psnr,raw_ssim,com_ssim' > $csvfile
for quality in `seq 5 5 90`; do
  r="raw_${quality}.jpg"
  c="com_${quality}.jpg"
  r_psnr=$(imagediff -m psnr $raw $r)
  c_psnr=$(imagediff -m psnr $raw $c)
  r_ssim=$(imagediff -m ssim $raw $r)
  c_ssim=$(imagediff -m ssim $raw $c)
  echo "${quality},${r_psnr},${c_psnr},${r_ssim},${c_ssim}" >> $csvfile
done
