# Makefile for Solitaire Fyne Application
# Project Variables
APP_NAME=iratesol
PACKAGE_ID=com.shanehowearth.iratesol
FYNE=fyne
SOURCE_DIR=cmd/fyne
DIST_DIR=dist

# Pathing
GUI_DIR=screen/gui

.PHONY: all android

default: native

android:
	fyne-cross android -app-id $(PACKAGE_ID) -icon cmd/fyne/Icon.png $(SOURCE_DIR)

deploy-android-test:
	adb install fyne-cross/dist/android/solitaire.apk

all: native asahi windows linux-x64 android

clean:
	@echo "Cleaning all build artifacts and temporary assets..."
	rm -rf $(DIST_DIR)
	rm -rf fyne-cross
	rm -rf $(TEMP_ASSET_DIR)

help:
	@echo "make android             - Build for Android"
	@echo "make deploy-android-test - Install to android device"
