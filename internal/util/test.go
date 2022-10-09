package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

//nolint:unused
func GetTestdataLocation() (string, error) {
	cmdOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get base directory: %v", err)
	}
	return fmt.Sprintf("%s/testdata", strings.TrimSpace(string(cmdOut))), nil
}

//nolint:deadcode,unused
func LoadTestConfig(t *testing.T, file string) string {
	loc, err := GetTestdataLocation()
	if err != nil {
		t.Fatalf("failed to get test data location")
	}

	data, err := os.ReadFile(fmt.Sprintf("%s/%s", loc, file))
	if err != nil {
		t.Fatalf("failed to load test config %s: %v", file, err.Error())
	}
	return string(data)
}

func LoadTestTemplateConfig(t *testing.T, file string, variables map[string]string) string {
	t1 := template.New("Config")
	t1, err := t1.Parse(LoadTestConfig(t, file))
	if err != nil {
		t.Fatalf("failed to parse config template %s: %v", file, err.Error())
	}

	var res bytes.Buffer
	if err = t1.Execute(&res, variables); err != nil {
		t.Fatalf("failed to insert variables into config template %s: %v", file, err.Error())
	}
	return res.String()
}
