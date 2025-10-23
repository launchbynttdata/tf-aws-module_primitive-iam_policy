#!/bin/bash

pyproject_path=$1
output_folder=$2
architecture=${3:-x86_64}
manylinux_version=${4:-manylinux_2_17}

# Wheels use "aarch64" rather than the "arm64" nomenclature that AWS uses
if [ $architecture = "arm64" ]; then architecture=aarch64; fi

echo "Creating Lambda Layer folder at $output_folder, building for $pyproject_path (platform: $manylinux_version_$architecture)"
rm -fr $output_folder
mkdir -p $output_folder
python3 -m pip install --platform $manylinux_version_$architecture --only-binary=:all: $pyproject_path -t $output_folder/python
