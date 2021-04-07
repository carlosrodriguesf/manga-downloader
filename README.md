### Tool to search mangas, download e generate MOBI files from the download content automatically

## Installation

If you don't want to use amazon kindle to read your manga, just skip to step 3

1) This application uses the [Kindle Comic Converter](https://github.com/ciromattia/kcc) to convert the downloaded
   content to a MOBI file. So if you want to use this type of file (to read on an amazon kindle for example), you need
   to install KCC. To do this just click [here]((https://github.com/ciromattia/kcc)) ant follow the official
   installation instructions.


2) It is important to note that kindle comic converter in turn uses kindlegen to generate `.mobi` files. Kindlegen has
   been discontinued by amazon and it is not possible to find official links to download it. Knowing this you can just
   click [here](https://github.com/carlosrodriguesf/manga-downloader/raw/master/vendor/kindlegen) and download binary.
   Once you have downloaded it just copy the file to the `/usr/bin` folder and give execute
   permission`chmod + x /usr/bin/kindlegen`. If you prefer you can look for another source to download on the internet
   too, just make sure to install kindlegen somehow or it will not be possible to generate the `.mobi` files.


3) To install this application is very simple, just
   click [here](https://github.com/carlosrodriguesf/manga-downloader/raw/master/cmd/mangadownloader) to download the
   binary and after downloading it open a terminal, navigate to the download folder and run it. After that, just follow
   the steps and in a few minutes you will have your manga downloaded and converted to `.mobi` (if you have followed
   steps 1 and 2). If you prefer to run globally you can simply copy the file to the `/usr/bin` folder. After that you
   can run it from the terminal from any directory.

Note: If you have the go environment set up on your machine you can just download the project and modify or compile it as is best for you. 

In case of doubt, send me an email: `contato@carlosrodrigues.dev.br`