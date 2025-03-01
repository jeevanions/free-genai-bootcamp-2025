#!/bin/bash

# Create and activate virtual environment using UV
echo "Creating virtual environment..."
uv venv

# Activate virtual environment
echo "Activating virtual environment..."
source .venv/bin/activate

# Install dependencies using UV
echo "Installing dependencies..."
uv pip install -e .

echo "Setup complete! Virtual environment is activated and dependencies are installed."
