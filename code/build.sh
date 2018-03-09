#!/usr/bin/bash
#
# dir structure
# /code
#     /bin
#     /log
#     /go   golang-binary
#     /src  src code
#     build.sh

WORKROOT=$(pwd)

go_root=$WORKROOT/go
go_path=$WORKROOT
dependence_path=$WORKROOT/src/mc2/trade/vendor

cd ${WORKROOT}

echo $WORKROOT

# unzip go environment
go_env="go1.8.3.linux-amd64.tar.gz"
if [ -d $go_root ]; then
	echo "go has been extracted before, remove it first"
	rm -rf $go_root
fi
tar -zxf third/$go_env
if [ $? -ne 0 ];
then
	echo "fail in extract go"
	exit 1
fi
echo "OK for extract go"

# unzip dependence package
if [ -d "${dependence_path}" ]; then
	echo "dependence package has been extracted before, remove it first"
	rm -rf $dependence_path
fi
mkdir $dependence_path
tar -zxf third/dependence.tar.gz -C ${dependence_path}
if [ $? -ne 0 ];
then
	echo "fail in extract dependence package"
	exit 1
fi

# prepare PATH, GOROOT and GOPATH
export PATH=$go_root/bin:$PATH
export GOROOT=$go_root
export GOPATH=$go_path


