package mkdir

import (
	"os"
	"path/filepath"

	"github.com/RomanIkonnikov93/niisva/internal/config"
)

func CreateStorageDir(cfg *config.Config) error {

	path := filepath.Clean(cfg.FileStoragePath)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}
