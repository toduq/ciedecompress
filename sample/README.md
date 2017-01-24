# 検証

# 圧縮率の検証

ciedecompressの0.1から2.0までのborder値でかけて、
それをPNGおよびJPEGで圧縮してみる。

それぞれの形式でのファイルサイズは以下のとおりであった。

![filesize](https://raw.githubusercontent.com/toduq/ciedecompress/master/sample/graph/b_filesize.png)

いずれの形式でもファイルサイズが減少していることが確認される。  
borderが1.0の値において、
PNGでは14％程度、JPEGでは3％程度減少していることから、  
PNGでは非常に効果的にciedecompressが働いていることがわかる。

これはPNGがもともと可逆圧縮であったためであると考えられる。

![psnr](https://raw.githubusercontent.com/toduq/ciedecompress/master/sample/graph/b_psnr.png)

PSNRの値はborderが1.0の部分においては、約52dB程度であり、  
目視でも劣化は感じられない。  
2.0においては顕著にブロックノイズが現れる。

SSIMはこの範囲において有意な値を示さなかったため、グラフは掲載しない。
必要であれば、`border/9_quality.csv`を確認していただきたい。

# JPEG圧縮との併用

非可逆圧縮であるJPEGと併用した場合について検証する。  
JPEGのqualityを変化させた時のファイルサイズが以下のとおりである。

![filesize](https://raw.githubusercontent.com/toduq/ciedecompress/master/sample/graph/q_filesize.png)

いずれの場合も僅かにciedecompressを使用した時のほうがファイルサイズは小さくなっている。
見た目の画質をほとんど変化させないまま、圧縮率を上げることができている。
ただし、PSNRは少しだけ悪化していることがわかる。

![psnr](https://raw.githubusercontent.com/toduq/ciedecompress/master/sample/graph/q_psnr.png)
![ssim](https://raw.githubusercontent.com/toduq/ciedecompress/master/sample/graph/q_ssim.png)

