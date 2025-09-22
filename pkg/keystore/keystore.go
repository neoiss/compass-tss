package keystore

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func DecryptKey(path string) (*keystore.Key, error) {
	// Make sure key exists before prompting password
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("key file not found: %s", path)
	}

	var pswd = GetPassword(fmt.Sprintf("Enter password for key %s:\n", path))

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
