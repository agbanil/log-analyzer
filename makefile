build:
	docker build -t gatsbytakehome .
run:
	docker run -e FILE_PATHS="file.txt" --rm --name gatsbytakehome gatsbytakehome
run-file:
	docker run -v $(file):/app/file.txt -e FILE_PATHS="file.txt" --rm --name gatsbytakehome gatsbytakehome
run-local:
	FILE_PATHS=file.txt go run .
test:
	go test -v