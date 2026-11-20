package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Sallle-1/task-3/config"
	"golang.org/x/text/encoding/charmap"
)

const outputDirMode = 0o755

var (
	errConfigPath  = errors.New("specify a YAML file with -config")
	errXMLEncoding = errors.New("unsupported XML encoding")
	errExtraXML    = errors.New("unexpected data after ValCurs")
	errEmptyRates  = errors.New("no Valute entries in XML")
	errCurrency    = errors.New("incomplete currency data")
)

type rate float64

func (value *rate) UnmarshalText(text []byte) error {
	decimal := strings.ReplaceAll(strings.TrimSpace(string(text)), ",", ".")

	number, err := strconv.ParseFloat(decimal, 64)
	if err != nil {
		return fmt.Errorf("parse currency value %q: %w", text, err)
	}

	*value = rate(number)

	return nil
}

type currency struct {
	NumCode  int    `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    *rate  `json:"value"     xml:"Value"`
}

func main() {
	path, err := readConfigPath()
	check(err)

	cfg, err := config.Load(path)
	check(err)

	rates, err := readCurrencies(cfg.InputFile)
	check(err)

	sort.SliceStable(rates, func(i, j int) bool {
		return *rates[i].Value > *rates[j].Value
	})

	check(saveJSON(cfg.OutputFile, rates))
}

func readConfigPath() (string, error) {
	path := flag.String("config", "", "path to the YAML file")
	flag.Parse()

	if *path == "" {
		return "", errConfigPath
	}

	return *path, nil
}

func readCurrencies(path string) ([]currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open XML %q: %w", path, err)
	}

	defer func() {
		check(file.Close())
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if !strings.EqualFold(charset, "windows-1251") {
			return nil, fmt.Errorf("%w %q", errXMLEncoding, charset)
		}

		return charmap.Windows1251.NewDecoder().Reader(input), nil
	}

	var document struct {
		XMLName    xml.Name   `xml:"ValCurs"`
		Currencies []currency `xml:"Valute"`
	}

	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("decode XML %q: %w", path, err)
	}

	var extra any

	if err := decoder.Decode(&extra); err == nil {
		return nil, errExtraXML
	} else if !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("decode XML %q: %w", path, err)
	}

	if len(document.Currencies) == 0 {
		return nil, errEmptyRates
	}

	for i := range document.Currencies {
		item := &document.Currencies[i]
		item.CharCode = strings.TrimSpace(item.CharCode)

		if item.Value == nil || (item.NumCode == 0) != (item.CharCode == "") {
			return nil, errCurrency
		}
	}

	return document.Currencies, nil
}

func saveJSON(path string, rates []currency) error {
	data, err := json.MarshalIndent(rates, "", "    ")
	if err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, outputDirMode); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	file, err := os.CreateTemp(dir, ".rates-*.json")
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}

	defer func() {
		if err := os.Remove(file.Name()); !errors.Is(err, os.ErrNotExist) {
			check(err)
		}
	}()

	_, writeErr := file.Write(data)
	closeErr := file.Close()

	if writeErr != nil {
		return fmt.Errorf("write JSON: %w", writeErr)
	}

	if closeErr != nil {
		return fmt.Errorf("close JSON file: %w", closeErr)
	}

	if err := os.Rename(file.Name(), path); err != nil {
		return fmt.Errorf("replace JSON file: %w", err)
	}

	return nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
