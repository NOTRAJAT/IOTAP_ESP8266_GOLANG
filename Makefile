buildTempl:
	@cd templ
	@templ generate 
	@cd ..
	@tailwindcss -i ./templ/css/output.css -o ./templ/css/build.css --minify


build:buildTempl
	@go build -o ./bin/build
	

run:build
	@./bin/build

test:
	@go test -v ./...	