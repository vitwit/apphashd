# Apphashd

A tool to quickly identify the root cause of apphash mismatches on Cosmos-based networks. Inspired by [apphash_calculator](https://gist.github.com/freak12techno/845a3061ed65295667c145c05ffd3b23).

## Features
- Compares two `application.db` directories from different nodes
- Calculates and compares module apphashes
- Identifies modules with differing hashes
- Uses [iaviewer](https://github.com/cosmos/iavl/tree/master/cmd/iaviewer) to pinpoint the exact data causing the mismatch

## Quick Start
1. **Download and run the script:**
   ```sh
   wget https://raw.githubusercontent.com/vitwit/apphashd/main/apphashd.sh && chmod 755 apphashd.sh
   ./apphashd.sh <path-to-application.db-node1> <path-to-application.db-node2>
   ```
2. **Results:**
   - A `hashes` folder is created inside `apphashd`.
   - Contains module hashes, IAVL trees, diffs, and decoded output showing the root cause of the mismatch.

## Usage Example
This repo includes test data from two nodes with an apphash mismatch.

### 1. Download test data
```sh
mkdir ~/node1 ~/node2
cd ~/node1
wget https://github.com/vitwit/apphashd/raw/main/testdata/node1/node1-application-db.tar.gz
 tar -xvf node1-application-db.tar.gz
cd ~/node2
wget https://github.com/vitwit/apphashd/raw/main/testdata/node2/node2-application-db.tar.gz
tar -xvf node2-application-db.tar.gz
```

### 2. Run the script
```sh
cd ~/
wget https://raw.githubusercontent.com/vitwit/apphashd/main/apphashd.sh && chmod 755 apphashd.sh
./apphashd.sh ~/node1/application.db ~/node2/application.db
```

### 3. Inspect the results
- Check the `diff-bank-decoded` file in `~/apphashd/hashes` for the root cause of the mismatch.

## How it Works
- Compares module hashes between two nodes
- Finds modules with differences
- Uses iaviewer to extract and diff IAVL trees
- Decodes the diff to reveal the exact data discrepancy

## More Information
- See [procedure.md](./procedure.md) for details on how the mismatch was simulated and debugged.
- Inspired by [freak12techno](https://github.com/freak12techno)'s [apphash_calculator](https://gist.github.com/freak12techno/845a3061ed65295667c145c05ffd3b23).
