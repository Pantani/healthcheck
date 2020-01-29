package fixtures

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

const (
	fixturesFolder = "configs"       // Folder contains the JSON fixtures
	fixturesFile   = "fixtures.json" // Fixtures for health check tests
)

func GeFixtures() (testFixtures Fixtures, err error) {
	err = geFixtures(fixturesFile, &testFixtures)
	return
}

func geFixtures(f string, r interface{}) error {
	b, err := getFile(f)
	if err != nil {
		return err
	}
	return json.Unmarshal(b[:], &r)
}

func getFile(file string) ([]byte, error) {
	golden := filepath.Join(fixturesFolder, file)
	return ioutil.ReadFile(golden)
}
