DUMP_URL = http://dumps.wikimedia.org/enwiki/latest/enwiki-latest-pages-articles.xml.bz2
DUMP_FILENAME = enwiki-latest-pages-articles.xml.bz2

#all: data

data/latest.xml.bz2:
	@mkdir -p data
	@wget -c $(DUMP_URL) -O data/$(DUMP_FILENAME)
	@ln -s $(DUMP_FILENAME) data/latest.xml.bz2
	@echo "Download completed!"

fmt: FORCE
	@go fmt *.go

run: FORCE
	@go run *.go

FORCE: