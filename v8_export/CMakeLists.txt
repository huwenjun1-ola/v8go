cmake_minimum_required(VERSION 3.22)
project(v8_export C CXX)

set(CMAKE_CXX_STANDARD 14)
add_definitions(-DCOMPILE_EXPORT)
include_directories(include)
add_definitions(-DCOMPILE_EXPORT)
add_definitions(-DV8_COMPRESS_POINTERS -DV8_31BIT_SMIS_ON_64BIT_ARCH)
add_library(v8_export STATIC V8InspectorClientImpl.h v8_export.h v8_export.cc)