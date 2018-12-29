#!/bin/bash

required_dir="bin conf logs src src/controller src/lib src/model src/configure src/logic src/proxy src/service src/utils"

function usage(){
    echo -e "Usage of generate.sh"
    echo -e "  -b string"
    echo -e "         required"
    echo -e "         project path, for example: /home/project/myapp"
    exit 1
}

root_path=""
package_name=""
while getopts "b:h*" Option
do
    case $Option in
        b) root_path=$OPTARG
        ;;
        h) usage
        ;;
        ?) 
        usage
        ;;
    esac
done

if [ "$root_path" == "" ]; then
    usage
fi

if [ "$package_name" == "" ]; then
    package_name=${root_path##*/}
fi

echo "您要创建的目录为: ${root_path}"
echo "您要创建的项目package名为: ${package_name}"
echo -e ""

read -p "确认创建[y/n]: " CONTINUE
if [ "n" == $CONTINUE ]; then
    exit 2
fi

here=`pwd`
transfer_dir=$here/transfer
mkdir -p $transfer_dir
cp -r ./src/* $transfer_dir

for f in ${transfer_dir}/*.go; do
    echo "replace package name:$package_name for ${f##*/}"
    sed -i "s/demo\//${package_name}\//g" $f
done

mkdir -p $root_path
cd $root_path

for d in $required_dir; do
    mkdir -p $d
done

cp $transfer_dir/serverctl $transfer_dir/make $transfer_dir/main.go $transfer_dir/README $root_path
cp -r $transfer_dir/*.toml $root_path/conf
cp $transfer_dir/DemoController.go $root_path/src/controller
cp $transfer_dir/DemoLogic.go $root_path/src/logic
cp $transfer_dir/Response.go $root_path/src/configure

cd $here

rm -rf $transfer_dir

echo "done"

exit 0
