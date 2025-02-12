#!/bin/sh
go build -o out/wall-collage main.go
cp wall-collage ~/.local/bin/wall-collage
chmod +x ~/.local/bin/wall-collage
