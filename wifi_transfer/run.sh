#!/bin/bash
# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# Check if uv is installed
if ! command -v uv &> /dev/null; then
    echo "Error: uv is not installed. Please install it first: https://github.com/astral-sh/uv"
    exit 1
fi

echo "Creating/Updating virtual environment..."
uv venv

echo "Installing dependencies..."
uv pip install -r requirements.txt -i https://mirrors.aliyun.com/pypi/simple/

echo "Starting WiFi Transfer Tool..."
.venv/bin/python main.py
