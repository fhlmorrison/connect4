TARGET=bin/app

run: build
	@./$(TARGET)

watch:
	air
	tailwindcss -i app.css -o public/app.css --watch

build: css-build
	go build -o $(TARGET) .

css-build: app.css
	tailwindcss -i app.css -o public/app.css --minify