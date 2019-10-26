#!/bin/bash
# This script was adapted from https://github.com/openfaas/cli.openfaas.com/blob/master/get.sh

version=$(curl -sI https://github.com/openfaas-incubator/ofc-bootstrap/releases/latest | grep Location | awk -F"/" '{ printf "%s", $NF }' | tr -d '\r')

if [ ! $version ]; then
    echo "Failed while attempting to install ofc-bootstrap. Please manually install:"
    echo ""
    echo "1. Open your web browser and go to https://github.com/openfaas-incubator/ofc-bootstrap/releases"
    echo "2. Download the latest release for your platform. Call it 'ofc-bootstrap'."
    echo "3. chmod +x ./ofc-bootstrap"
    echo "4. mv ./ofc-bootstrap /usr/local/bin"
    exit 1
fi

hasCli() {

    has=$(which ofc-bootstrap)

    if [ "$?" = "0" ]; then
        echo
        echo "You already have the ofc-bootstrap cli!"
        export n=1
        sleep $n
    fi

    hasCurl=$(which curl)
    if [ "$?" = "1" ]; then
        echo "You need curl to use this script."
        exit 1
    fi
}

getPackage() {
    uname=$(uname)
    userid=$(id -u)

    suffix=""
    case $uname in
    "Darwin")
    suffix="-darwin"
    ;;
    "Linux")
        arch=$(uname -m)
        echo $arch
        case $arch in
        "aarch64")
        suffix="-arm64"
        ;;
        esac
        case $arch in
        "armv6l" | "armv7l")
        suffix="-armhf"
        ;;
        esac
    ;;
    esac

    targetFile="/tmp/ofc-bootstrap$suffix"

    if [ "$userid" != "0" ]; then
        targetFile="$(pwd)/ofc-bootstrap$suffix"
    fi

    if [ -e $targetFile ]; then
        rm $targetFile
    fi

    url=https://github.com/openfaas-incubator/ofc-bootstrap/releases/download/$version/ofc-bootstrap$suffix
    echo "Downloading package $url as $targetFile"

    curl -sSLf $url --output $targetFile

    if [ "$?" = "0" ]; then

    chmod +x $targetFile

    echo "Download complete."

        if [ "$userid" != "0" ]; then

            echo
            echo "========================================================="
            echo "==    As the script was run as a non-root user the     =="
            echo "==    following commands may need to be run manually   =="
            echo "========================================================="
            echo
            echo "  sudo cp ofc-bootstrap$suffix /usr/local/bin/ofc-bootstrap"
            echo

        else

            echo
            echo "Running as root - Attempting to move ofc-bootstrap to /usr/local/bin"

            mv $targetFile /usr/local/bin/ofc-bootstrap

            if [ "$?" = "0" ]; then
                echo "New version of ofc-bootstrap installed to /usr/local/bin"
            fi

            if [ -e $targetFile ]; then
                rm $targetFile
            fi

           ofc-bootstrap -version
        fi
    fi
}

hasCli
getPackage
