# Apphash Mismatch Debugging Guide

This guide explains a practical, step-by-step procedure to identify the root cause of an apphash mismatch on Cosmos-based networks. The goal is to pinpoint the exact source of non-determinism that causes a chain halt.

---

## 1. Setup: Simulating an Apphash Mismatch
- Two-node devnet using [Juno](https://github.com/CosmosContracts/juno/) binaries
- Validator 1: v22.0.0
- Validator 2: v22.0.1
- A transaction triggers an apphash mismatch, halting the chain

---

## 2. Step-by-Step Debugging Procedure

### Step 1: Identify the Offending Module
- Use an apphash calculator (e.g., [apphash_calculator](https://gist.github.com/freak12techno/845a3061ed65295667c145c05ffd3b23)) on both nodes:
  ```sh
  go run main.go <path-to-application.db>
  ```
- Example output (Validator 1):
  ```
  got commitInfo with 34 stores
   index 0: store name 08-wasm, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
   ...
   index 3: store name bank, hash 8699243f40720dc1c10f65457f04f7be3e88305e1d571917ab38807dc22d733f
   ...
  hash: 6a84014dfe03092c897c17b979e47b2374c7df6f0b6a993abee1fee2ee146ce7
  ```
- Example output (Validator 2):
  ```
  got commitInfo with 34 stores
   ...
   index 3: store name bank, hash eea82414b653f802a9e1127add39e4aac569e72b4b36beb07036a88b3ae726c8
   ...
  hash: 94800fc8dc32db79388b5033d40d091b21b870865889a9d0b9e7563db93c44b9
  ```
- Focus on the module(s) with differing hashes (e.g., `bank`)

### Step 2: Inspect the Offending Module's IAVL Tree
- Use [iaviewer](https://github.com/cosmos/iavl/tree/master/cmd/iaviewer) to dump the IAVL tree shape for the module:
  ```sh
  iaviewer shape <path-to-application.db> s/k:<module>/ > shape-<node>.txt
  # Example for bank: s/k:bank/
  ```
- Example output (Validator 1):
  ```
      *4 00666163746F72792F6A756E6F31686B6B706A6B6D6779676A33777579397473637474616E74766B793239787379747374356D7A2F74657374
    -3 DF426B01E303332BC40C69CBE92BEBA05E926EB73D2A765857E453B87F2C906A
    ...
  ```
- Example output (Validator 2):
  ```
      *4 00666163746F72792F6A756E6F31686B6B706A6B6D6779676A33777579397473637474616E74766B793239787379747374356D7A2F74657374
    -3 B2B690FE501397C91C4ED9B17A04F45F480EF71D98F401CDAF10BF792753E3A3
    ...
  ```
- Compare the outputs:
  ```sh
  diff shape-1.txt shape-2.txt
  ```
- Example diff output:
  ```
  1a2,3
  >         *4 00666163746F72792F6A756E6F31686B6B706A6B6D6779676A33777579397473637474616E74766B793239787379747374356D7A2F74657374
  >       -3 B2B690FE501397C91C4ED9B17A04F45F480EF71D98F401CDAF10BF792753E3A3
  3c5
  <       -3 DF426B01E303332BC40C69CBE92BEBA05E926EB73D2A765857E453B87F2C906A
  ---
  >     -2 36A2B0B9EFB8576635CC592EB8A2767FA1DA8D0D41185E75CCD86BA5EF8EEA9F
  ...
  ```
- Look for differences in the tree structure or keys

### Step 3: Decode the Difference
- Convert any differing hex keys to ASCII to identify the problematic data:
  ```sh
  echo "00666163746F72792F6A756E6F31686B6B706A6B6D6779676A33777579397473637474616E74766B793239787379747374356D7A2F74657374" | xxd -r -p
  # Output:
  factory/juno1hkkpjkmgygj3wuy9tscttantvky29xsytst5mz/test
  ```
- The decoded key/value points to the root cause of the mismatch

---

## 3. Example: How the Apphash Mismatch Was Triggered
- Created a new denom using the `tokenfactory` module:
  ```sh
  junod tx tokenfactory create-denom test --from <address> --chain-id test
  ```
- Minted tokens for the new denom:
  ```sh
  junod tx tokenfactory mint 1000000factory/<address>/test --from <address> --chain-id test
  ```
- This transaction triggered the mismatch due to a code difference between the two binary versions

---

## 4. Tips & Notes
- Module names for iaviewer must be in the format `s/k:<module>/` (e.g., `s/k:staking/`)
- If you see `Error reading data: version does not exist`, check your module name format
- Compare binary versions for recent changes if you find a mismatch

---

## References
- [Juno v22.0.0...v22.0.1 diff](https://github.com/CosmosContracts/juno/compare/v22.0.0...v22.0.1)
- [apphash_calculator gist](https://gist.github.com/freak12techno/845a3061ed65295667c145c05ffd3b23)
