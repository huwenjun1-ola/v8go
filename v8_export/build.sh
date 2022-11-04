#!/bin/sh
export CC=clang
export CXX=clang++
if [ -e CMakeCache.txt ]; then
  rm CMakeCache.txt
fi
mkdir -p build_cmake
cd build_cmake
cmake -DCMAKE_BUILD_TYPE=Release ../
make