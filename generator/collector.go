package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/Unknwon/goconfig"
)

// Collector defines the interface for collecting config and translation data
type Collector interface {
	Collect(string) (Data, error)
	AvailableISOs() []string
}

// ErrNoConfig is the error to be returned if there are no config files
var ErrNoConfig = fmt.Errorf("no config files")

type Data struct {
	I18N   map[string]Map
	Config map[string]Map
}

func (d Data) String() string {
	b, err := json.Marshal(d)
	if err != nil {
		log.Printf(err.Error())
		return ""
	}
	return string(b)
}

type Map map[string]string

func (d Map) String() string {
	b, err := json.Marshal(d)
	if err != nil {
		log.Printf(err.Error())
		return ""
	}
	return string(b)
}

func NewCollector(configFolder, i18nFolder string) Collector {
	return SimpleCollector{ConfigFolder: configFolder, I18NFolder: i18nFolder}
}

type SimpleCollector struct {
	ConfigFolder string
	I18NFolder   string
}

func (c SimpleCollector) Collect(iso string) (Data, error) {
	data := Data{}
	translations, err := c.getTranslations(iso)
	if err != nil {
		return data, err
	}
	data.I18N = translations

	configs, err := c.getConfigurations(iso)
	if err != nil {
		return data, err
	}
	data.Config = configs

	return data, nil
}

func (c SimpleCollector) AvailableISOs() []string {
	tmp := map[string]struct{}{}
	for _, iso := range c.availableConfigISOs() {
		tmp[iso] = struct{}{}
	}
	for _, iso := range c.availableTranslationISOs() {
		tmp[iso] = struct{}{}
	}
	isos := []string{}
	for iso := range tmp {
		isos = append(isos, iso)
	}
	sort.Strings(isos)
	return isos
}

func (c SimpleCollector) availableTranslationISOs() []string {
	isos := []string{}
	files, err := ioutil.ReadDir(c.I18NFolder)
	if err != nil {
		log.Printf("looking for available translation ISO codes: %s", err.Error())
		return isos
	}

	for _, f := range files {
		if !f.IsDir() && strings.Contains(f.Name(), ".ini") {
			isos = append(isos, strings.Trim(f.Name(), ".ini"))
		}
	}
	return isos
}

func (c SimpleCollector) availableConfigISOs() []string {
	isos := []string{}
	files, err := ioutil.ReadDir(c.ConfigFolder)
	if err != nil {
		log.Printf("looking for available config ISO codes: %s", err.Error())
		return isos
	}

	for _, f := range files {
		if f.IsDir() && f.Name() != "global" {
			isos = append(isos, f.Name())
		}
	}
	return isos
}

func (c SimpleCollector) getConfigurations(iso string) (map[string]Map, error) {
	configFiles := []string{}

	for _, p := range []string{
		fmt.Sprintf("%s/global", c.ConfigFolder),
		fmt.Sprintf("%s/%s", c.ConfigFolder, iso),
	} {
		files, err := ioutil.ReadDir(p)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		for _, f := range files {
			configFiles = append(configFiles, fmt.Sprintf("%s/%s", p, f.Name()))
		}
	}

	result := map[string]Map{}

	extra := []string{}
	switch len(configFiles) {
	case 0:
		return result, ErrNoConfig
	case 1:
	default:
		extra = configFiles[1:]
	}
	configs, err := goconfig.LoadConfigFile(configFiles[0], extra...)
	if err != nil {
		return result, err
	}

	for _, section := range configs.GetSectionList() {
		configSection, err := configs.GetSection(section)
		if err != nil {
			return result, err
		}
		log.Println("loaded config section", section, configSection)
		result[section] = configSection
	}

	return result, nil
}

func (c SimpleCollector) getTranslations(iso string) (map[string]Map, error) {
	result := map[string]Map{}
	translations, err := goconfig.LoadConfigFile(fmt.Sprintf("%s/%s.ini", c.I18NFolder, iso))
	if err != nil {
		return result, err
	}

	for _, section := range translations.GetSectionList() {
		configSection, err := translations.GetSection(section)
		if err != nil {
			return result, err
		}
		log.Println("loaded translation section", section, "with", len(configSection), "translations")
		result[section] = configSection
	}

	return result, nil
}
