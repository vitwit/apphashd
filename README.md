# Apphashd

This is a framework for identifying the root cause of an appphash mismatch on cosmos based networks inspired by [apphash_calculator](https://gist.github.com/freak12techno/845a3061ed65295667c145c05ffd3b23) written by [freak12techno](https://github.com/freak12techno). 

The `apphashd` binary takes two inputs as arguments, absolute/relative paths to the `application.db` of nodes with opposing apphashes. It calculates and compares the apphashes of all the modules present in the db. It generates an array with the name of the modules which have differing hashes. 

This repo contains a script called `apphashd.sh` which is an end to end script which builds the `apphashd` binary, executes it and uses [iaviewer](https://github.com/cosmos/iavl/tree/master/cmd/iaviewer) to identify the message_type/data that caused the apphash to deviate on both the nodes. 

## Usage
Download and run the script
```
wget https://raw.githubusercontent.com/vitwit/apphashd/main/apphashd.sh && chmod 755 apphashd.sh
./apphashd.sh <path-to-application.db-node1> <path-to-application.db-node2>
```
The script creates a folder called `hashes` inside the `apphashd` folder which contains the module hashes of both the nodes, iavl trees of the differing modules, diff of iavl trees and the ascii decoded output of the diff which contains the root cause of apphash mismatch. 

### Example usage
This repo has a folder called `testdata` which contains the application.db of two nodes which were halted due to an apphash mismatch. We will use that to test the script.

#### Download the testdata
```
mkdir ~/node1 ~/node2
cd ~/node1
wget https://github.com/vitwit/apphashd/raw/main/testdata/node1/node1-application-db.tar.gz
tar -xvf node1-application-db.tar.gz
cd ~/node2
wget https://github.com/vitwit/apphashd/raw/main/testdata/node2/node2-application-db.tar.gz
tar -xvf node2-application-db.tar.gz
```

#### Download and execute the script
```
cd ~/
wget https://raw.githubusercontent.com/vitwit/apphashd/main/apphashd.sh && chmod 755 apphashd.sh
./apphashd.sh ~/node1/application.db ~/node2/application.db
```

The script generated a file called `diff-bank-decoded` inside the `~/apphashd/hashes` folder. We can check the root cause of the apphash mismatch by inspecting this file.
```
factory/juno1hkkpjkmgygj3wuy9tscttantvky29xsytst5mz/test
	Ebzf>�����`Up6ɳk\vW�]Z��;Y
��6��N�.�v�X7m�;��&�\�ӻ�6Q
�X?R�%�����J��ğ5��'�І�qc�xv
�!9�E�ErB�l&��E��d^;��2���
�@(x�Y�Bk5/RJ�jstake
woGϬ�U��ds*�;�+e�,ʣ^�
Q@�Sn��������n���)�H���3�
�@(x�Y�Bk5/RJ�jstake
��Mۗ�幗�t��pW�cU�N,��A�=jG�O
?,vD�t�S���Ѽ/lBq�+��cR��>n
��[h"%p�\0��ke���factory/juno1hkkpjkmgygj3wuy9tscttantvky29xsytst5mz/test
s;�ܝ"s�]�����<u�+"��rs{n�
��[h"%p�\0��ke���stake
�D����f�V/�<4@����$�6?r�I���
factory/juno1hkkpjkmgygj3wuy9tscttantvky29xsytst5mz/test��[h"%p�\0��ke���
���D��?",
���я9/�4٥	!����KU�/[}�_
��[h"%p�\0��ke���stake
����v���������Q���*���)b��
```

Refer to this [doc](./procedure.md) to know how the apphash mismatch was simulated in a two node devnet and how [freak12techno's](https://github.com/freak12techno) [apphash_calculator](https://gist.github.com/freak12techno/845a3061ed65295667c145c05ffd3b23) was used to debug it originally which served as the inspiration for `apphashd` 
