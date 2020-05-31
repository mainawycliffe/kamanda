#!/usr/bin/env bash
set -xe

# determine os
case $(uname -s) in
    Darwin) platform='darwin';;
    *) platform='linux';;
esac

# determine architecture
case $(uname -m) in
    "x86_64") arch="amd64";; 
    *) arch="386";;
esac


if [ $# -eq 0 ]; then
    # download the latest executables for kamanda 
	kamanda_asset_path=$(
		command curl -sSf https://github.com/mainawycliffe/kamanda/releases |
			command grep -o "/mainawycliffe/kamanda/releases/download/.*/kamanda_.*_${platform}_${arch}.zip" |
			command head -n 1
	)
    if [ ! "$kamanda_asset_path" ]; then exit 1; fi
    downloadUrl="https://github.com/${kamanda_asset_path}"
else
    # download the specific requested version for kamanda
    downloadUrl="https://github.com/mainawycliffe/kamanda/releases/download/v${1}/kamanda_${1}_${platform}_${arch}.zip"
fi

binDir="/usr/local/bin"
filename="kamanda"

curl --fail --location --progress-bar "${downloadUrl}" --output "./${filename}.zip"

# extract executables
unzip kamanda.zip

# add to use executable dir, no permissions required
mv ${filename} ${binDir}

# delete the extra files
echo "Clean up the dowloads after installation"
rm "./${filename}.zip"

echo "Kamanda was installed successfully. You can find documentation here https://kamanda.dev"
