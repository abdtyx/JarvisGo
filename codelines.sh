#!/bin/bash
find . -type f -name "*.go" -exec cat {} + | wc -l
