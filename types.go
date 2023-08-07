package main

type MetaField struct {
	Name               string   `json:"name"`
	Type               int      `json:"type"`
	Notes              string   `json:"notes,omitempty"`
	InternalReferences []string `json:"internal_references" yaml:"internal_references"`
}

type MetaFile struct {
	Name   string      `json:"name"`
	Type   int         `json:"type"`
	Notes  string      `json:"notes,omitempty"`
	Fields []MetaField `json:"fields"`
}

type Column struct {
	Name string    `yaml:"name,omitempty"`
	Meta MetaField `yaml:"meta"`
}

type Model struct {
	Name    string    `yaml:"name"`
	Meta    MetaField `yaml:"meta"`
	Columns []Column  `yaml:"columns"`
}

type DbtYaml struct {
	Models []Model `yaml:"models"`
}
