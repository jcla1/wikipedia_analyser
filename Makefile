DUMP_URL = http://dumps.wikimedia.org/enwiki/20130708/enwiki-20130708-pages-articles.xml.bz2
DUMP_FILENAME = enwiki-20130708-pages-articles.xml.bz2
DUMP_MD5 = ce66b6b08514ddfc5c2296e3bbbd42fa

#all: data

data/latest.xml.bz2:
	@mkdir -p data
	#@wget -c $(DUMP_URL) -O data/$(DUMP_FILENAME)
	@ln -s $(DUMP_FILENAME) data/latest.xml.bz2
	@echo "Download completed!"