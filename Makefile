.PHONY: android
android:
	GOOS=android GOARCH=arm64 go build go_static\cmd\mp4server   

.PHONY: build
build:
	go build go_static\cmd\mp4server

.DEFAULT_GOAL := build