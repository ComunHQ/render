#!/bin/bash
find . -name "*.generated.yaml" -type f -delete
go run ../main.go render.yaml
