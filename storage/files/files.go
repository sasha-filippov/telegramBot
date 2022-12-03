package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"telegramBot/lib/e"
	"telegramBot/storage"
	"time"
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	fPath := filepath.Join(s.basePath, page.UserName)
	// check if a folder exist
	if _, err := os.Stat(fPath); os.IsNotExist(err) {
		if err := os.MkdirAll(fPath, os.ModePerm); err != nil {
			msg := fmt.Sprintf("value of fileperm is %w", os.ModePerm)
			return errors.New(msg)
		}
	}

	fName, err := filename(page)
	if err != nil {
		return errors.New("couldn't hash")
	}
	fPath = filepath.Join(fPath, fName)
	file, err := os.Create(fPath)
	if err != nil {
		return errors.New("couldn't create a file")
	}
	defer func() { _ = file.Close() }()
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return errors.New("Couldn't gobbing!")
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.New("no folder created")
	}
	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}
	//0- max-1
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]
	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	filename, err := filename(p)
	if err != nil {
		return e.Wrap("can't remove a page", err)
	}
	path := filepath.Join(s.basePath, p.UserName, filename)
	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove a file %s", path)
		return e.Wrap(msg, err)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	filename, err := filename(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}
	path := filepath.Join(s.basePath, p.UserName, filename)
	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)
		return false, e.Wrap(msg, err)

	}
	return true, nil

}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil
}

func filename(p *storage.Page) (string, error) {
	return p.Hash()
}
