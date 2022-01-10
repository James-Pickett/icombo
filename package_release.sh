#!/bin/bash

set -e

git_branch=$(git rev-parse --abbrev-ref HEAD)
oses=(
    "linux" 
    "windows" 
    "darwin"
    )
arches=(
    "386" 
    "amd64" 
    "arm64"
    )

artifactDir="./artifacts"
rm -rf "${artifactDir}"
mkdir -p "${artifactDir}"

generate_release_artifacts(){
    local os=$1
    local arch=$2

    # set up dir
    binDirName="icombo_${git_branch}_${os}_${arch}"

    echo "building binary package for ${binDirName}"

    binDirPath="${artifactDir}/${binDirName}"
    mkdir -p "${binDirPath}"
    
    # build binary
    binName="icombo"
    if [ "$os" == "windows" ]
    then
        binName="${binName}.exe"
    fi
    
    GOOS=$os GOARCH=$arch go build -o "${binDirPath}/${binName}"

    # package binary
    if [ "$os" == "windows" ]
    then
        zip -rjD "${artifactDir}/${binDirName}.zip" "${binDirPath}/"
    else
        tarPath="${artifactDir}/${binDirName}.tar.gz"
        tar -C "${binDirPath}" -czvf "${tarPath}" "${binName}"
    fi

    echo ""
    echo "building example package for ${binDirName}"
    
    # add example file
    cp -r "./example/"* "${binDirPath}"

    # remove outputs so users can see it work 
    rm -rf "${binDirPath}/output_images"

    # package example project
    if [ "$os" == "windows" ]
    then
        (cd "${binDirPath}" && zip -r "../../artifacts/${binDirName}_example_project.zip" .)
    else
        tarPath="${artifactDir}/${binDirName}_example_project.tar.gz"
        tar -C "${binDirPath}" -czvf "${tarPath}" "."
    fi

    rm -rf "${binDirPath}"

    echo ""
}

for os in "${oses[@]}"
do
    for arch in "${arches[@]}"
    do
        if [ "$os" == "darwin" ] && [ "$arch" == "386" ]
        then
            continue
        fi
        generate_release_artifacts $os $arch
    done
done