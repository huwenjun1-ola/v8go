cd "$(dirname "$0")"
./build.sh
cd ./build_cmake
cp ./libv8_export.* ../../lib/linux_x86_64/
cd ../
rm -rf ./build_cmake