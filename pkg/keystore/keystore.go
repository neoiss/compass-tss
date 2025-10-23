package keystore

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

var keystorePwd = "KEYSTORE_PASSWORD"

func DecryptKey(path string) (*keystore.Key, error) {
	// Make sure key exists before prompting password
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("key file not found: %s", path)
	}

	var pswd = []byte(os.Getenv(keystorePwd))
	if len(pswd) == 0 {
		pswd = GetPassword("Enter keystore password: ")
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read keyFile failed, err:%s", err)
	}
	ret, err := keystore.DecryptKey(file, string(pswd))
	if err != nil {
		return nil, fmt.Errorf("DecryptKey failed, err:%s", err)
	}

	return ret, nil
}
