APP_NAME=mp4server
CMD_PATH=go_static/cmd/mp4server
BUILD_FLAGS=-ldflags="-s -w"
OUTPUT=bin/$(APP_NAME)

install: build
	cp $(OUTPUT) /usr/bin/$(APP_NAME)

.PHONY: android
android:
	GOOS=android GOARCH=arm64 go build $(BUILD_FLAGS) $(CMD_PATH) 

.PHONY: build
build:
	go build $(BUILD_FLAGS) -o $(OUTPUT) $(CMD_PATH) 

clean:
	rm $(OUTPUT)

.DEFAULT_GOAL := build