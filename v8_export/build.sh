#!/bin/sh

if [ -e CMakeCache.txt ]; then
  rm CMakeCache.txt
fi
mkdir -p build_cmake
cmake -S ./ -B ./build_cmake
cd build_cmake
make