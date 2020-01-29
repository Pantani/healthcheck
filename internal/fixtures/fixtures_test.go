package fixtures

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeFixtures(t *testing.T) {
	_, err := GeFixtures()
	if err == nil {
		t.Errorf("GeFixtures() error = %v, wantErr %v", err, true)
	}
}

func Test_geFixtures(t *testing.T) {
	type args struct {
		f string
		r interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    Fixtures
		wantErr bool
	}{
		{"test success", args{f: "testdata/" + fixturesFile, r: &Fixtures{}}, Fixtures{{"ethereum", "https://eth1.trezor.io", []Test{{"block_height", "api", "blockbook.bestHeight", "GET", "5s", "lastValue <= newValue", map[string]interface{}{}}}}}, false},
		{"test error,,1", args{f: "blockatlas", r: nil}, nil, true},
		{"test error,,2", args{f: "blockatlas", r: &Fixtures{}}, nil, true},
		{"test error,,3", args{f: "../../main.go", r: &Fixtures{}}, Fixtures{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := geFixtures(tt.args.f, tt.args.r)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.EqualValues(t, &tt.want, tt.args.r)
		})
	}
}

func Test_getFile(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    []byte
		wantErr bool
	}{

		{"test success", "testdata/" + fixturesFile, []byte{91, 10, 32, 32, 123, 10, 32, 32, 32, 32, 34, 110, 97, 109, 101, 34, 58, 32, 34, 101, 116, 104, 101, 114, 101, 117, 109, 34, 44, 10, 32, 32, 32, 32, 34, 104, 111, 115, 116, 34, 58, 32, 34, 104, 116, 116, 112, 115, 58, 47, 47, 101, 116, 104, 49, 46, 116, 114, 101, 122, 111, 114, 46, 105, 111, 34, 44, 10, 32, 32, 32, 32, 34, 116, 101, 115, 116, 115, 34, 58, 32, 91, 10, 32, 32, 32, 32, 32, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 110, 97, 109, 101, 34, 58, 32, 34, 98, 108, 111, 99, 107, 95, 104, 101, 105, 103, 104, 116, 34, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 109, 101, 116, 104, 111, 100, 34, 58, 32, 34, 71, 69, 84, 34, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 117, 114, 108, 95, 112, 97, 116, 104, 34, 58, 32, 34, 97, 112, 105, 34, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 106, 115, 111, 110, 95, 112, 97, 116, 104, 34, 58, 32, 34, 98, 108, 111, 99, 107, 98, 111, 111, 107, 46, 98, 101, 115, 116, 72, 101, 105, 103, 104, 116, 34, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 101, 120, 112, 114, 101, 115, 115, 105, 111, 110, 34, 58, 32, 34, 108, 97, 115, 116, 86, 97, 108, 117, 101, 32, 60, 61, 32, 110, 101, 119, 86, 97, 108, 117, 101, 34, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 98, 111, 100, 121, 34, 58, 32, 123, 125, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 117, 112, 100, 97, 116, 101, 95, 116, 105, 109, 101, 34, 58, 32, 34, 53, 115, 34, 10, 32, 32, 32, 32, 32, 32, 125, 10, 32, 32, 32, 32, 93, 10, 32, 32, 125, 10, 93}, false},
		{"test error,3", "../../main.go", []byte{112, 97, 99, 107, 97, 103, 101, 32, 109, 97, 105, 110, 10, 10, 105, 109, 112, 111, 114, 116, 32, 40, 10, 9, 34, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 116, 114, 117, 115, 116, 119, 97, 108, 108, 101, 116, 47, 104, 101, 97, 108, 116, 104, 99, 104, 101, 99, 107, 47, 99, 109, 100, 34, 10, 41, 10, 10, 47, 47, 32, 64, 116, 105, 116, 108, 101, 32, 84, 114, 117, 115, 116, 87, 97, 108, 108, 101, 116, 32, 77, 101, 116, 114, 105, 99, 115, 32, 67, 111, 108, 108, 101, 99, 116, 111, 114, 10, 47, 47, 32, 64, 118, 101, 114, 115, 105, 111, 110, 32, 49, 46, 48, 10, 47, 47, 32, 64, 100, 101, 115, 99, 114, 105, 112, 116, 105, 111, 110, 32, 67, 111, 108, 108, 101, 99, 116, 32, 109, 101, 116, 114, 105, 99, 115, 32, 97, 110, 100, 32, 102, 111, 114, 109, 97, 116, 32, 102, 111, 114, 32, 80, 114, 111, 109, 101, 116, 104, 101, 117, 115, 10, 10, 47, 47, 32, 64, 99, 111, 110, 116, 97, 99, 116, 46, 110, 97, 109, 101, 32, 84, 114, 117, 115, 116, 32, 87, 97, 108, 108, 101, 116, 10, 47, 47, 32, 64, 99, 111, 110, 116, 97, 99, 116, 46, 117, 114, 108, 32, 104, 116, 116, 112, 115, 58, 47, 47, 116, 46, 109, 101, 47, 119, 97, 108, 108, 101, 99, 111, 114, 101, 10, 10, 47, 47, 32, 64, 108, 105, 99, 101, 110, 115, 101, 46, 110, 97, 109, 101, 32, 77, 73, 84, 32, 76, 105, 99, 101, 110, 115, 101, 10, 47, 47, 32, 64, 108, 105, 99, 101, 110, 115, 101, 46, 117, 114, 108, 32, 104, 116, 116, 112, 115, 58, 47, 47, 114, 97, 119, 46, 103, 105, 116, 104, 117, 98, 117, 115, 101, 114, 99, 111, 110, 116, 101, 110, 116, 46, 99, 111, 109, 47, 116, 114, 117, 115, 116, 119, 97, 108, 108, 101, 116, 47, 114, 101, 100, 101, 109, 112, 116, 105, 111, 110, 47, 109, 97, 115, 116, 101, 114, 47, 76, 73, 67, 69, 78, 83, 69, 10, 102, 117, 110, 99, 32, 109, 97, 105, 110, 40, 41, 32, 123, 10, 9, 99, 109, 100, 46, 69, 120, 101, 99, 117, 116, 101, 40, 41, 10, 125, 10}, false},
		{"test error,1", "blockatlas", nil, true},
		{"test error,2", "blockatlas", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFile(tt.file)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.EqualValues(t, tt.want, got)
		})
	}
}