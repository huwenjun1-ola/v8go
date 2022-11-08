cmake_minimum_required(VERSION 3.10)
project(v8_test C CXX)


set(CMAKE_CXX_STANDARD 14)

set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -nostdlib++ -stdlib=libc++ -static")
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -nostdlib++ -stdlib=libc++ -static")
set(CMAKE_EXE_LINKER_FLAGS "-fuse-ld=lld")
include_directories(${CMAKE_SOURCE_DIR}/../)
link_directories(${CMAKE_SOURCE_DIR}/../lib/linux_x86_64/)
add_definitions(-DCOMPILE_EXPORT)
add_definitions(-DV8_COMPRESS_POINTERS -DV8_31BIT_SMIS_ON_64BIT_ARCH)
add_executable(test main.c)
target_link_libraries(test libv8_export.a libv8.a libv8_libbase.a libv8_libplatform.a  atomic dl pthread rt c++ c++abi m)