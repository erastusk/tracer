package cache

import (
	"encoding/json"
	"fmt"
	"os"

	"github/erastusk/tracer/internal/types"
	"github/erastusk/tracer/internal/utils"
)

const (
	dirFile = ".tracer"
)

type Cache struct {
	fileName string
}

func NewCache(filename string) *Cache {
	return &Cache{
		fileName: filename,
	}
}

func (c *Cache) LoadCache() types.PromptOptions {
	var ans types.PromptOptions
	f, err := c.openFile()
	if err != nil {
		utils.Green.Println(err)
	}
	json.NewDecoder(f).Decode(&ans)
	f.Close()
	return ans
}

func (c *Cache) openFile() (*os.File, error) {
	file := fmt.Sprintf(".tracer/%s.json", c.fileName)
	if _, err := os.Stat(file); err == nil {
		return os.OpenFile(file, os.O_RDWR, 0o644)
	}
	os.MkdirAll(dirFile, os.ModePerm)
	os.Chmod(dirFile, 0o755)
	f, err := os.Create(file)
	if err != nil {
		utils.Green.Println("error creating file:", file)
		return nil, err
	}
	return f, err
}

func (c *Cache) SaveCache(ans types.PromptOptions) {
	f, _ := c.openFile()
	json.NewEncoder(f).Encode(ans)
	f.Close()
}
