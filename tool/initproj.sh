#!/bin/sh
###################################
#       初始化项目目录脚本          #
###################################
read -p "Enter Project name:" name
if [ ! -n "$name" ]; then
    echo "项目名不能为空"
    exit 1
fi

if [[ ! "$name" =~ ^[A-Z][A-Za-z0-9]*$ ]]; then
    echo "项目首字母必须大写"
    exit 1
fi

#检测目录是否已经存在
newName="PG${name}"
if [ -d "$newName" ]; then
    echo "项目已经存在"
    exit 1
fi
#项目目录结构
array_dir_name[0]=${newName}/app/controllers
array_dir_name[1]=${newName}/app/entity
array_dir_name[2]=${newName}/app/proto
array_dir_name[3]=${newName}/app/service
array_dir_name[4]=${newName}/conf
array_dir_name[5]=${newName}/logs

for i in ${array_dir_name[@]}; do
    if [ ! -d "$i" ]; then
        mkdir -p "$i"
    fi
done

