package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/store/rootmulti"
	dbm "github.com/cosmos/cosmos-db"
)

type ModuleHash struct {
	StoreName string `json:"store_name"`
	Hash      string `json:"hash"`
}

func getModuleHashes(dbPath string) ([]ModuleHash, error) {
	dbDir := filepath.Dir(dbPath)
	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, dbDir)
	if err != nil {
		return nil, err
	}

	multistoreRaw := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	multistore := multistoreRaw.(*rootmulti.Store)
	version := multistore.LatestVersion()
	commitInfo, err := multistore.GetCommitInfo(version)
	if err != nil {
		return nil, err
	}

	var moduleHashes []ModuleHash
	for _, storeInfo := range commitInfo.StoreInfos {
		moduleHashes = append(moduleHashes, ModuleHash{
			StoreName: storeInfo.Name,
			Hash:      fmt.Sprintf("%x", storeInfo.GetHash()),
		})
	}

	// Sort the moduleHashes slice alphabetically by StoreName
	sort.Slice(moduleHashes, func(i, j int) bool {
		return moduleHashes[i].StoreName < moduleHashes[j].StoreName
	})

	return moduleHashes, nil
}

func saveHashesToFile(moduleHashes []ModuleHash, fileName string) error {
	jsonOutput, err := json.MarshalIndent(moduleHashes, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, jsonOutput, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Hashes saved to %s\n", fileName)
	fmt.Println(string(jsonOutput))
	return nil
}

func main() {
	if len(os.Args) != 3 {
		panic("please provide exactly 2 paths to application.db")
	}

	dbPath1 := os.Args[1]
	dbPath2 := os.Args[2]

	moduleHashes1, err := getModuleHashes(dbPath1)
	if err != nil {
		panic(err)
	}

	moduleHashes2, err := getModuleHashes(dbPath2)
	if err != nil {
		panic(err)
	}

	// Create the hashes directory in $PWD if it doesn't exist
	outputDir := filepath.Join(".", "hashes")
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Directory hashes already exists. Overwriting files...")
	}

	// Save the hashes in node-1-hashes.json and node-2-hashes.json
	file1 := filepath.Join(outputDir, "node-1-hashes.json")
	file2 := filepath.Join(outputDir, "node-2-hashes.json")

	err = saveHashesToFile(moduleHashes1, file1)
	if err != nil {
		panic(err)
	}

	err = saveHashesToFile(moduleHashes2, file2)
	if err != nil {
		panic(err)
	}

	// Compare the module hashes between the two databases
	differingModules := []string{}
	for i, hash1 := range moduleHashes1 {
		hash2 := moduleHashes2[i]
		if hash1.StoreName == hash2.StoreName && hash1.Hash != hash2.Hash {
			fmt.Printf("Differing module: %s\n", hash1.StoreName)
			fmt.Printf("DB1 Hash: %s\n", hash1.Hash)
			fmt.Printf("DB2 Hash: %s\n", hash2.Hash)
			differingModules = append(differingModules, hash1.StoreName)
		}
	}

	// Export differing modules as an array in a new file
	if len(differingModules) > 0 {
		envFilePath := filepath.Join(outputDir, "modules.env")
		envFile, err := os.Create(envFilePath)
		if err != nil {
			panic(err)
		}
		defer envFile.Close()

		// Format the differing modules as a bash array
		envVarValue := fmt.Sprintf("DIFFERING_MODULES=(%s)\n", strings.Join(differingModules, " "))

		_, err = envFile.WriteString(envVarValue)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Differing modules exported to %s as an array: %s\n", envFilePath, envVarValue)
	} else {
		fmt.Println("No differing modules found.")
	}
}
