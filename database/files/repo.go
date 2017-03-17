package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/StageAutoControl/controller/cntl"
	"github.com/StageAutoControl/controller/cntl/dmx"
	"gopkg.in/yaml.v2"
)

type fileData struct {
	SetLists        []*cntl.SetList    `json:"setLists" yaml:"setLists"`
	Songs           []*cntl.Song       `json:"songs" yaml:"songs"`
	DmxScenes       []*dmx.Scene       `json:"dmxScenes" yaml:"dmxScenes"`
	DmxPresets      []*dmx.Preset      `json:"dmxPresets" yaml:"dmxPresets"`
	DmxAnimations   []*dmx.Animation   `json:"dmxAnimations" yaml:"dmxAnimations"`
	DmxDevices      []*dmx.Device      `json:"dmxDevices" yaml:"dmxDevices"`
	DmxDeviceGroups []*dmx.DeviceGroup `json:"dmxDeviceGroups" yaml:"dmxDeviceGroups"`
}

// Repository is a file repository
type Repository struct {
	dataDir string
}

// New crates a new file repository and returns it.
func New(dataDir string) *Repository {
	return &Repository{
		dataDir: dataDir,
	}
}

// Load implements cntl.Loader and loads the data from filesystem
func (r *Repository) Load() (*cntl.DataStore, error) {
	store := cntl.NewStore()
	return r.readDir(store, r.dataDir)
}

func (r *Repository) readDir(data *cntl.DataStore, dir string) (*cntl.DataStore, error) {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range fs {
		path := filepath.Join(dir, f.Name())

		if f.IsDir() {
			data, err = r.readDir(data, path)
			if err != nil {
				return nil, err
			}

			continue
		}

		ext := filepath.Ext(path)
		switch ext {
		case ".yml", ".yaml":
			data, err = r.readYAMLFile(data, path)
			if err != nil {
				return nil, err
			}

			break

		case ".json":
			data, err = r.readJSONFile(data, path)
			if err != nil {
				return nil, err
			}

			break

		default:
			return nil, fmt.Errorf("Unable to load file %q. No loader for file extension %q known.", path, ext)
		}
	}

	return data, nil
}

func (r *Repository) readJSONFile(data *cntl.DataStore, path string) (*cntl.DataStore, error) {
	var fileData fileData

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading file %q: %v", path, err)
	}

	err = json.Unmarshal(b, &fileData)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse content of %q: %v", path, err)
	}

	return r.mergeData(data, &fileData), nil
}

func (r *Repository) readYAMLFile(data *cntl.DataStore, path string) (*cntl.DataStore, error) {
	var fileData fileData

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading file %q: %v", path, err)
	}

	err = yaml.Unmarshal(b, &fileData)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse content of %q: %v", path, err)
	}

	return r.mergeData(data, &fileData), nil
}

func (r *Repository) mergeData(data *cntl.DataStore, fileData *fileData) *cntl.DataStore {
	for _, sl := range fileData.SetLists {
		data.SetLists[sl.ID] = sl
	}

	for _, s := range fileData.Songs {
		data.Songs[s.ID] = s
	}

	for _, d := range fileData.DmxDevices {
		data.DmxDevices[d.ID] = d
	}

	for _, dg := range fileData.DmxDeviceGroups {
		data.DmxDeviceGroups[dg.ID] = dg
	}

	for _, p := range fileData.DmxPresets {
		data.DmxPresets[p.ID] = p
	}

	for _, sc := range fileData.DmxScenes {
		data.DmxScenes[sc.ID] = sc
	}

	return data
}
