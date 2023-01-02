#!/bin/bash
set -e

# version and keys are supplied as arguments
version="$1"
rc=`echo $version | awk -F - '{print $2}'`
keyfile="$2"
pubkeyfile="$3" # optional

if [[ -z $version || -z $keyfile ]]; then
	echo "Usage: $0 VERSION KEYFILE"
	exit 1
fi
if [[ -z $pubkeyfile ]]; then
	echo "Warning: no public keyfile supplied. Binaries will not be verified."
fi

# check for keyfile before proceeding
if [ ! -f $keyfile ]; then
    echo "Key file not found: $keyfile"
    exit 1
fi

# import key from $keyfile to gpg keys
gpg --import $keyfile

# setup build-time vars
ldflags="-s -w -X 'github.com/EvilRedHorse/pubaccess-node/build.GitRevision=`git rev-parse --short HEAD`' -X 'github.com/EvilRedHorse/pubaccess-node/build.BuildTime=`date`' -X 'github.com/EvilRedHorse/pubaccess-node/build.ReleaseTag=${rc}'"

for arch in amd64 arm; do
	for os in darwin linux windows; do

		# We don't need ARM versions for windows or mac (just linux)
		if [ "$arch" == "arm" ]; then
			if [ "$os" == "windows" ] || [ "$os" == "darwin" ] || [ "$os" == "freebsd" ]; then
				continue
			fi
		fi

		echo Packaging ${os} ${arch}...

		# create workspace
		folder=release/pubaccess-node-$version-$os-$arch
		# move older builds to new folder, appending move time in Nanoseconds since UNIX epoch to folder name
		mv -f $folder $folder$(printf "_mv_$(date +%s%N)_ns")
		mkdir -p $folder
		# compile and sign binaries
		for pkg in spc spd; do
			bin=$pkg
			if [ "$os" == "windows" ]; then
				bin=${pkg}.exe
			fi
			CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -a -tags 'netgo' -ldflags="$ldflags" -o $folder/$bin ./cmd/$pkg

		done

		# add other artifacts ### changing after v1.6.1
		# cp -r LICENSE README.md $folder
		# zip
		(
			cd release
			zip -rq pubaccess-node-$version-$os-$arch.zip pubaccess-node-$version-$os-$arch
			# sign zip release
			gpg --armour --output pubaccess-node-$version-$os-$arch.zip.asc --detach-sig pubaccess-node-$version-$os-$arch.zip
		)
	done
done
