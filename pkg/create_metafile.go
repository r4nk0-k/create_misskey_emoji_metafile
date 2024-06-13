package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type JsonFile struct {
	MetaVersion int     `json:"metaVersion"`
	Host        string  `json:"host"`
	ExportedAt  string  `json:"exportedAt"`
	Emojis      []Emoji `json:"emojis"`
}

type Emoji struct {
	FileName   string      `json:"fileName"`
	Downloaded bool        `json:"downloaded"`
	Emoji      EmojiDetail `json:"emoji"`
}

type EmojiDetail struct {
	ID                                      string   `json:"id"`
	UpdatedAt                               string   `json:"updatedAt"`
	Name                                    string   `json:"name"`
	Host                                    string   `json:"host"`
	Category                                string   `json:"category"`
	OriginalUrl                             string   `json:"originalUrl"`
	PublicUrl                               string   `json:"publicUrl"`
	Uri                                     string   `json:"uri"`
	Type                                    string   `json:"type"`
	Aliases                                 []string `json:"aliases"`
	License                                 string   `json:"license"`
	LocalOnly                               bool     `json:"localOnly"`
	IsSensitive                             bool     `json:"isSensitive"`
	RoleIdsThatCanBeUsedThisEmojiAsReaction []string `json:"roleIdsThatCanBeUsedThisEmojiAsReaction"`
}

type ConfigYamlSchema struct {
	Host           string                `yaml:"host"`
	EmojiParameter ConfigYamlEmojiSchema `yaml:"emojiParameter"`
}

type ConfigYamlEmojiSchema struct {
	License     string `yaml:"license"`
	IsSensitive bool   `yaml:"isSensitive`
	LocalOnly   bool   `yaml:"localonly"`
	Category    string `yaml:"category"`
}

const (
	ConfigYaml = "./cfg/config.yaml"
)

func main() {
	dn := os.Args[1]
	t := time.Now().Format("2006-01-02T04:05:06Z")

	yf, err := os.Open(ConfigYaml)
	if err != nil {
		log.Fatal(err)
	}
	defer yf.Close()

	y, err := io.ReadAll(yf)
	if err != nil {
		log.Fatal(err)
	}

	var conf ConfigYamlSchema
	if err := yaml.Unmarshal(y, &conf); err != nil {
		log.Fatal(err)
	}

	j := &JsonFile{
		MetaVersion: 2,
		Host:        conf.Host,
		ExportedAt:  t,
	}

	emojis := []Emoji{}
	_ = filepath.WalkDir(dn, func(path string, d fs.DirEntry, err error) error {
		if ext := filepath.Ext(path); ext == ".png" || ext == ".PNG" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".GIF" {
			pat := regexp.MustCompile(`.*\(\d+\).*`) //(1)とか(2)がついているファイルは除外
			if pat.MatchString(path) {
				fmt.Println("duplicate file:", path)
				return nil
			}

			fn := strings.ReplaceAll(filepath.Base(path), "-", "_")
			os.Rename(path, filepath.Dir(path)+"/"+fn)

			ed := EmojiDetail{
				Name:        fn[:len(filepath.Base(path))-len(filepath.Ext(path))],
				Category:    conf.EmojiParameter.Category,
				LocalOnly:   conf.EmojiParameter.LocalOnly,
				IsSensitive: conf.EmojiParameter.IsSensitive,
				License:     conf.EmojiParameter.License,
				Type:        "image/webp",
				Aliases:     []string{},
			}

			e := Emoji{
				FileName:   fn,
				Downloaded: true,
				Emoji:      ed,
			}

			emojis = append(emojis, e)
		}
		return nil
	})

	j.Emojis = emojis
	if d, err := json.Marshal(j); err != nil {
		log.Fatal(err)
	} else {
		var buf bytes.Buffer
		err := json.Indent(&buf, []byte(d), "", "  ")
		if err != nil {
			panic(err)
		}
		f, err := os.Create(dn + "/meta.json")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if _, err := f.Write(buf.Bytes()); err != nil {
			log.Fatal(err)
		}
	}
}
