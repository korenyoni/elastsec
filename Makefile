default:
	go build

clean:
	@rm -f secbeat

zip_linux_amd64:
	zip elastsec_$(git describe $(git rev-list --tags --max-count=1))_linux_amd64 elastsec
