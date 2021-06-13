.PHONY: android
android:
	GOOS=android GOARCH=arm64 go build -o go_static_android

.PHONY: linux
linux:
	go build

.DEFAULT_GOAL := linux