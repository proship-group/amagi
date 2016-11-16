package storage

import (
	"fmt"
	"os"
	"strings"
	"time"

	utils "github.com/b-eee/amagi"
	vault "github.com/hashicorp/vault/api"
)

var (
	// VaultClient vault main client
	VaultClient *vault.Client

	// SecretBackend secret backend prefix
	SecretBackend = "secret"
)

type (
	request map[string]interface{}
)

// StartVault start vault connections
func StartVault() error {
	config := vault.DefaultConfig()
	utils.Info(fmt.Sprintf("StartVault connecting to -->> %v", config.Address))

	client, err := vault.NewClient(config)
	if err != nil {
		utils.Error(fmt.Sprintf("error StartVault %v", err))
		return err
	}

	client.SetToken(os.Getenv("ROOT_TOKEN"))

	VaultClient = client

	fmt.Println("StartVault connected!")
	return nil
}

// LogicalClient logical client
func LogicalClient() *vault.Logical {
	return VaultClient.Logical()
}

// PathName path name builders
func PathName(paths ...string) string {
	return strings.Join(paths, "/")
}

// SplitRedisKey split redis key for paths
func SplitRedisKey(backend, key string) []string {
	b := []string{backend}
	b = append(b, strings.Split(key, ":")...)

	return b
}

// VWrite vault write
func VWrite(path string, data request) error {
	s := time.Now()
	// for k, v := range data {
	// 	fmt.Println("k", k, "v", v, "===============")
	// }

	if _, err := LogicalClient().Write(path, data); err != nil {
		utils.Error(fmt.Sprintf("error VWrite %v", err))
		return err
	}

	utils.Info(fmt.Sprintf("VWrite took: %v path=%v", time.Since(s), path))
	return nil
}

// VRead vault read path
func VRead(paths ...string) (*vault.Secret, error) {
	s := time.Now()
	path := PathName(paths...)
	secret, err := LogicalClient().Read(path)
	if err != nil {
		utils.Error(fmt.Sprintf("error on VRead %v", err))
		return secret, err
	}

	for k := range secret.Data {
		v, err := LogicalClient().Unwrap(secret.Data[k].(string))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(v, "===xx")
	}

	utils.Info(fmt.Sprintf("VRead took: %v path=%v", time.Since(s), path))
	return secret, nil
}
