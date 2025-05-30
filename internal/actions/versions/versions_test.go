package versions

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"golang.org/x/mod/semver"
)

func TestCVersions(t *testing.T) {
	b := []byte(`{
		"cgood": {
			"versions" : {
				"1.3": ["v0.1.0", "v0.1.1"],
				"1.3.1": ["v1.1.0"]
			}
		}
	}`)
	path := "ttt.json"
	err := os.WriteFile(path, []byte(b), 0755)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(path)

	v := Read(path)
	goVersion := v.GoVersions("cgood")
	semver.Sort(goVersion)

	if !reflect.DeepEqual(goVersion, []string{"v0.1.0", "v0.1.1", "v1.1.0"}) {
		t.Errorf("unexpected goversion: want: %v got: %v", []string{"v0.1.0", "v0.1.1", "v1.1.0"}, goVersion)
	}

	cVersion := v.CVersions("cgood")
	semver.Sort(cVersion)

	if !reflect.DeepEqual(cVersion, []string{"v1.3.0", "v1.3.1"}) {
		t.Errorf("unexpected cversion: want: %v got: %v", []string{"v1.3.0", "v1.3.1"}, cVersion)
	}

	if v.LatestGoVersionForCVersion("cgood", "1.3") != "v0.1.1" {
		t.Errorf("unexpected latest Go version: want: %v got: %v", "v0.1.1", v.LatestGoVersionForCVersion("cgood", "1.3"))
	}

	if v.SearchBySemVer("cgood", "v1.3.0") != "1.3" {
		t.Errorf("unexpected search by semver result: want: %v got: %v", "1.3", v.SearchBySemVer("cgood", "v1.3.0"))
	}

	if v.SearchBySemVer("agood", "v1.3.0") != "" {
		t.Errorf("unexpected search by semver result: want: %v got: %v", "", v.SearchBySemVer("cgood", "v1.3.0"))
	}
}

func TestLatestVersion(t *testing.T) {
	b := []byte(`{
    "cgood": {
        "versions" : {
				"1.3": ["v0.1.0", "v0.1.1"],
				"1.3.1": ["v1.1.0"]
		}
    }
}`)
	path := "ttt.json"
	err := os.WriteFile(path, []byte(b), 0755)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(path)

	v := Read(path)

	if v.LatestGoVersion("cgood") != "v1.1.0" {
		t.Errorf("unexpected latest version: want: v1.1.0 got: %s", v.LatestGoVersion("cgood"))
	}
}

func TestAppend(t *testing.T) {
	v := Read("llpkgstore.json")
	defer os.Remove("llpkgstore.json")

	v.Write("cjson", "1.7.18", "v1.0.0")
	v.Write("cjson", "1.7.19", "v1.0.2")

	v = Read("llpkgstore.json")
	//defer os.Remove("llpkgstore.json")

	v.Write("cjson", "1.7.18", "v1.0.1")
	v.Write("libxml", "1.45.1.4", "v1.0.0")

	v = Read("llpkgstore.json")
	v.Write("libxml", "1.45.1.5", "v1.0.1")

	b, _ := os.ReadFile("llpkgstore.json")

	if !bytes.Equal(b, []byte(`{"cjson":{"versions":{"1.7.18":["v1.0.0","v1.0.1"],"1.7.19":["v1.0.2"]}},"libxml":{"versions":{"1.45.1.4":["v1.0.0"],"1.45.1.5":["v1.0.1"]}}}`)) {
		t.Error("unexpected append result")
	}
}
