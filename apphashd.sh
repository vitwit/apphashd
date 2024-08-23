#!/bin/bash

function color() {
    # Usage: color "31;5" "string"
    # Some valid values for color:
    # - 5 blink, 1 strong, 4 underlined
    # - fg: 31 red,  32 green, 33 yellow, 34 blue, 35 purple, 36 cyan, 37 white
    # - bg: 40 black, 41 red, 44 blue, 45 purple
    printf '\033[%sm%s\033[0m\n' "$@"
}

## Function for checking if a given command exists on the system.
command_exists () {
    type "$1" &> /dev/null ;
}

## Checking if the required arguments are provided to the script. 
if [ $# -eq 2 ]; then
    color "33" "Paths to the dbs provided"
else
    color "31" "The script requires 2 arguments"
    exit
fi

## Updating system packages to install prerequisites for building the apphashd binary  
color "33" "Updating system packages"
sleep 3
sudo apt-get update
sudo apt-get -y upgrade
sudo apt-get -y dist-upgrade
sudo apt-get install build-essential -y

color "33" "Checking for go installation on the system"
sleep 3
if command_exists go ; then
  color "32" "Golang is already installed."
  sleep 3
else
  color "33" "Golang is not installed. Proceeding with go installation"
  sleep 3
  cd /tmp
  wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz
  tar -xvf go1.22.2.linux-amd64.tar.gz
  sudo mv go /usr/local
  mkdir -p ~/go/src/github.com
  mkdir ~/go/bin
  export GOROOT=/usr/local/go
  export GOPATH=$HOME/go
  export GOBIN=$GOPATH/bin
  export PATH=$PATH:/usr/local/go/bin:$GOBIN

  echo "" >> ~/.bashrc
  echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
  echo 'export GOPATH=$HOME/go' >> ~/.bashrc
  echo 'export GOBIN=$GOPATH/bin' >> ~/.bashrc
  echo 'export PATH=$PATH:/usr/local/go/bin:$GOBIN' >> ~/.bashrc
fi

## Checking if apphashd is already built and installed on the system
if command_exists apphashd ; then
  color "32" "apphashd is already installed."
  cd apphashd
  sleep 3
else
## Downloading and building apphashd binary
  color "33" "Downloading repo and building the apphashd binary"
  sleep 3
  cd ~/
  git clone https://github.com/vitwit/apphashd.git && cd apphashd 
  go build -o apphashd
  sudo cp apphashd /usr/local/bin/
fi

## Executing the binary to find differing hashes.
DB1=$1
DB2=$2
apphashd $1 $2
source ./hashes/modules.env
DIFFERING_MODULES_COUNT=${#DIFFERING_MODULES[@]}
if (( $DIFFERING_MODULES_COUNT > 4 )); then
  color "31" "NUmber of modules with differing hashes is high. Please ensure you've provided the paths to the correct databases with the same height. If the correct databases are provided you might be using an incompatible version of iavl. If you still want to continue then press ENTER. Press CTRL+C to exit the script"
  read 
fi

## Downloading and building iaviewer tool
rm -rf iavl
git clone https://github.com/cosmos/iavl.git && cd iavl/cmd/iaviewer
go build -o iaviewer
sudo cp iaviewer /usr/local/bin/
cd ../../..

for module in $DIFFERING_MODULES; do
color "33" "Printing the iavl tree of $module of db1"
sleep 3
touch hashes/node-1-$module-shape
iaviewer shape $DB1 s/k:$module/ > hashes/node-1-$module-shape
color "33" "Printing the iavl tree of $module of db2"
sleep 3
touch hashes/node-2-$module-shape
iaviewer shape $DB2 s/k:$module/ > hashes/node-2-$module-shape
diff hashes/node-1-$module-shape hashes/node-2-$module-shape > hashes/diff-$module-shape

## ASCII decoding of hex string
while IFS= read -r line; do
  while [[ "$line" =~ ([*-][0-9]+\ )([0-9A-F]+) ]]; do
    hex="${BASH_REMATCH[2]}"
    ascii=$(echo "$hex" | xxd -r -p | tr -d '\000')
    echo "$ascii" >> hashes/diff-$module-decoded
    line="${line#*${BASH_REMATCH[2]}}"
  done
done < hashes/diff-$module-shape

done
