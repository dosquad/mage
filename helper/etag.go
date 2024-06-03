package helper

import (
	"errors"
	"io"
	"os"

	"github.com/na4ma4/go-permbits"
	"gopkg.in/yaml.v3"
)

type ETagItem struct {
	parent *ETag
	Key    string `yaml:"name"`
	Value  string `yaml:"value"`
}

func (e *ETagItem) Save() error {
	e.parent.Set(e.Key, e.Value)
	return e.parent.Save()
}

type ETag []ETagItem

func ETagLoadConfig() (ETag, error) {
	etag := ETag{}

	var f *os.File
	{
		var err error
		f, err = os.Open(MustGetArtifactPath(".etag.yml"))
		if errors.Is(err, os.ErrNotExist) {
			return etag, nil
		} else if err != nil {
			return etag, err
		}
	}
	defer f.Close()

	{
		if err := yaml.NewDecoder(f).Decode(&etag); err != nil {
			if errors.Is(err, io.EOF) {
				return etag, nil
			}

			return etag, err
		}
	}

	return etag, nil
}

func (e *ETag) Set(key, value string) {
	for idx := range *e {
		if (*e)[idx].Key == key {
			(*e)[idx].Value = value
			return
		}
	}

	(*e) = append((*e), ETagItem{
		Key:   key,
		Value: value,
	})
}

func (e ETag) GetItem(key string) *ETagItem {
	for _, v := range e {
		if v.Key == key {
			v.parent = &e
			return &v
		}
	}

	return &ETagItem{
		parent: &e,
		Key:    key,
	}
}

func (e ETag) GetRelative(path string) string {
	after, ok := GetRelativePath(path)
	if ok {
		return e.Get(after)
	}

	return path
}

func (e ETag) Get(key string) string {
	for _, v := range e {
		if v.Key == key {
			return v.Value
		}
	}

	return ""
}

func (e ETag) Save() error {
	MustMakeDir(MustGetArtifactPath(), permbits.MustString("ug=rwx,o=rx"))

	var f *os.File
	{
		var err error
		f, err = os.Create(MustGetArtifactPath(".etag.yml"))
		if err != nil {
			return err
		}
	}
	defer f.Close()

	if err := yaml.NewEncoder(f).Encode(e); err != nil {
		return err
	}

	return nil
}
