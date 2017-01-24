import matplotlib
import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
import datetime as dt
plt.style.use('ggplot')
font = {'family' : 'TakaoPGothic'}
matplotlib.rc('font', **font)

import os
import glob
import re

# PNGが圧縮できていることを示す
borders = []
png = []
jpg = []
png_max = os.stat("border/b0.0.png").st_size
jpg_max = os.stat("border/b0.0.jpg").st_size
for filename in sorted(glob.glob('border/b*.png')):
	borders.append(float(re.findall(r"border/b(.*)\.png", filename)[0]))
	png.append(os.stat(filename).st_size / png_max * 100)
for filename in sorted(glob.glob('border/b*.jpg')):
	jpg.append(os.stat(filename).st_size / jpg_max * 100)

df = pd.DataFrame({
	'png': png,
	'jpg': jpg,
}, index=borders)

df.plot(y=['png', 'jpg'])
plt.title('borderの値とファイルサイズの関係', size=16)
plt.xlabel("border")
plt.ylabel("ファイルサイズ [%]")
plt.legend([
	"PNG形式(100％=%s[KB])" % str(png_max//1000),
	"JPEG形式(100％=%s[KB])" % str(jpg_max//1000)
], loc="lower left")
plt.savefig("graph/b_filesize.png")

df = pd.read_csv('border/9_quality.csv', index_col='border')
df.plot(y=['psnr'])
plt.title('borderの値によるPSNRの変化', size=16)
plt.xlabel("border")
plt.ylabel("Peak signal-to-noise ratio [dB]")
plt.legend(["PSNR"], loc="upper left")
plt.savefig("graph/b_psnr.png")

df = pd.read_csv('border/9_quality.csv', index_col='border')
df.plot(y=['ssim'])
plt.title('borderの値によるSSIMの変化', size=16)
plt.xlabel("border")
plt.ylabel("Structural similarity")
plt.legend(["SSIM"], loc="upper left")
plt.savefig("graph/b_ssim.png")

# JPEGにおいて同一のqualityでも圧縮をかけたほうがサイズが小さくなることを示す
qualities = []
withs = []
withouts = []
for quality in range(5, 90, 5):
	qualities.append(quality)
	withs.append(float(os.stat("quality/com_%s.jpg" % quality).st_size) / 1000)
	withouts.append(float(os.stat("quality/raw_%s.jpg" % quality).st_size) / 1000)

df = pd.DataFrame({
	'with': withs,
	'without': withouts
}, index=qualities)

df.plot(y=['without', 'with'])
plt.title('ciedecompressありと無しでのファイルサイズの差', size=16)
plt.xlabel("jpeg quality")
plt.ylabel("ファイルサイズ [KB]")
plt.legend(["ciedecompressなし", "ciedecompressあり"], loc="upper left")
plt.savefig("graph/q_filesize.png")

df = pd.read_csv('quality/9_quality.csv', index_col='quality')
df.plot(y=['raw_psnr', 'com_psnr'])
plt.title('ciedecompressによるPSNRの差', size=16)
plt.xlabel("jpeg quality")
plt.ylabel("Peak signal-to-noise ratio [dB]")
plt.legend(["ciedecompressなし", "ciedecompressあり"], loc="upper left")
plt.savefig("graph/q_psnr.png")

df = pd.read_csv('quality/9_quality.csv', index_col='quality')
df.plot(y=['raw_ssim', 'com_ssim'])
plt.title('ciedecompressによるSSIMの差', size=16)
plt.xlabel("jpeg quality")
plt.ylabel("Structural similarity")
plt.legend(["ciedecompressなし", "ciedecompressあり"], loc="upper left")
plt.savefig("graph/q_ssim.png")
