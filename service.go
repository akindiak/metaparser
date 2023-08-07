package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v2"
)

type MetaParser interface {
	Parse(string, *string, *[]string) error
	GetPaths(string) ([]string, error)
	CreateMetaFile(string, *Model) error
}

type YmlParser struct{}

func (yp *YmlParser) Parse(fpath string, system *string, tables *[]string) error {
	// Read the file
	ymlData, err := os.ReadFile(fpath)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	var data DbtYaml

	err = yaml.Unmarshal(ymlData, &data)
	if err != nil {
		fmt.Println(err)
	}

	for _, m := range data.Models {
		if contains(tables, m.Name) {
			err = yp.CreateMetaFile(*system, &m)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
		}
	}
	return nil
}

func (yp *YmlParser) GetPaths(dirPath string) ([]string, error) {
	var filePaths []string

	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".yml" {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return filePaths, nil
}

func (yp *YmlParser) CreateMetaFile(system string, m *Model) error {
	dirPath := fmt.Sprintf("metafiles/%s", system)
	_ = os.Mkdir(dirPath, os.ModePerm)

	filePath := path.Join(dirPath, fmt.Sprintf("%s.json", m.Name))
	fmt.Println(filePath)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	var fields []MetaField

	for _, c := range m.Columns {
		if !reflect.ValueOf(c.Meta).IsZero() {
			fields = append(fields, MetaField{
				Name:               c.Name,
				Type:               c.Meta.Type,
				Notes:              c.Meta.Notes,
				InternalReferences: c.Meta.InternalReferences,
			})
		}
	}

	metaTable := &MetaFile{
		Name:   m.Name,
		Type:   m.Meta.Type,
		Notes:  m.Meta.Notes,
		Fields: fields,
	}

	encoder := json.NewEncoder(f)
	if err = encoder.Encode(&metaTable); err != nil {
		return err
	}

	return nil
}

func contains(s *[]string, str string) bool {
	for _, v := range *s {
		if v == str {
			return true
		}
	}

	return false
}
