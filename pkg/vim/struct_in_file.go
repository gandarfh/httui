package vim

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/google/uuid"
)

type Preview struct {
	data interface{}
	file string
}

func NewPreview(data interface{}) Preview {
	file := filepath.Join(os.TempDir(), uuid.New().String()+".json")

	return Preview{data, file}
}

func (v *Preview) create_tmp_file() error {
	file, err := os.Create(v.file)
	if err != nil {
		return errors.BadRequest("Error when try create the file.\n", err.Error())
	}

	text, err := json.MarshalIndent(v.data, "", "\t")
	if err != nil {
		return errors.BadRequest("Error when try marshal resource data to json file.\n", err.Error())
	}

	file.Write([]byte(text))

	return nil
}

func (v *Preview) Open() error {
	if err := v.create_tmp_file(); err != nil {
		return err
	}

	cmd := exec.Command("lvim", v.file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()

	return nil
}

func (v *Preview) Execute(decode interface{}) error {
	file, err := os.ReadFile(v.file)
	if err != nil {
		return errors.BadRequest("Error when try read the file.\n", err.Error())
	}

	if err := json.Unmarshal(file, decode); err != nil {
		return errors.BadRequest("Error when try marshal to repository.\n", err.Error())
	}

	return nil
}

func (v *Preview) Close() error {
	if err := os.Remove(v.file); err != nil {
		return errors.BadRequest("Error when try delete the tmp file.\n", err.Error())
	}

	return nil
}
