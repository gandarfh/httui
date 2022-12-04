package terminal

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type Finish struct {
	Err      error
	Category string
	Preview  *Preview
}

type Preview struct {
	Data interface{}
	File string
}

func NewPreview(data interface{}) *Preview {
	file := filepath.Join(os.TempDir(), uuid.New().String()+".json")

	return &Preview{data, file}
}

func (p *Preview) create_tmp_file() error {
	file, err := os.Create(p.File)
	if err != nil {
	}

	text, err := json.MarshalIndent(p.Data, "", "\t")
	if err != nil {
	}

	file.Write([]byte(text))

	return nil
}

func (p *Preview) OpenVim(category string) tea.Cmd {
	p.create_tmp_file()
	cmd := exec.Command("nvim", p.File)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return Finish{err, category, p}
	})
}

func (p *Preview) Execute(decode interface{}) error {
	file, _ := os.ReadFile(p.File)

	json.Unmarshal(file, decode)

	return nil
}

func (p *Preview) Close() error {
	if err := os.Remove(p.File); err != nil {
	}

	return nil
}
