LATEST_TAG := $(shell git describe $(shell git rev-list --tags --max-count=1))

default:
	go build

clean:
	@rm -f secbeat

zip_linux_amd64:
	zip elastsec_$(LATEST_TAG)_linux_amd64.zip elastsec && md5sum elastsec_$(LATEST_TAG)_linux_amd64.zip > elastsec_$(LATEST_TAG)_linux_amd64.md5sum
