package v1

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	v1pb "github.com/bytebase/bytebase/proto/generated-go/v1"
)

func TestGenerateEtag(t *testing.T) {
	tests := []struct {
		statement string
		want      string
	}{
		{
			statement: "",
			want:      "da39a3ee5e6b4b0d3255bfef95601890afd80709",
		},
		{
			statement: "test",
			want:      "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
		},
		{
			statement: "CREATE TABLE test;",
			want:      "ecbe3800ad8d24592e2dc963ae63f96c608723db",
		},
	}

	for _, test := range tests {
		got := GenerateEtag([]byte(test.statement))
		require.Equal(t, test.want, got)
	}
}

type transformTest struct {
	Engine   v1pb.Engine
	Schema   string
	Metadata *v1pb.DatabaseMetadata
}

func TestTransformSchemaString(t *testing.T) {
	const (
		record = false
	)
	var (
		filepath = "testdata/schema.yaml"
	)

	a := require.New(t)
	yamlFile, err := os.Open(filepath)
	a.NoError(err)

	tests := []transformTest{}
	byteValue, err := io.ReadAll(yamlFile)
	a.NoError(yamlFile.Close())
	a.NoError(err)
	a.NoError(yaml.Unmarshal(byteValue, &tests))

	for i, t := range tests {
		result, err := transformSchemaStringToDatabaseMetadata(t.Engine, t.Schema)
		a.NoError(err)
		if record {
			tests[i].Metadata = result
		} else {
			a.Equal(t.Metadata, result)
		}
	}

	if record {
		byteValue, err := yaml.Marshal(tests)
		a.NoError(err)
		err = os.WriteFile(filepath, byteValue, 0644)
		a.NoError(err)
	}
}

type designTest struct {
	Engine   v1pb.Engine
	Baseline string
	Target   *v1pb.DatabaseMetadata
	Result   string
}

func TestGetDesignSchema(t *testing.T) {
	const (
		record = false
	)
	var (
		filepath = "testdata/design.yaml"
	)

	a := require.New(t)
	yamlFile, err := os.Open(filepath)
	a.NoError(err)

	tests := []designTest{}
	byteValue, err := io.ReadAll(yamlFile)
	a.NoError(yamlFile.Close())
	a.NoError(err)
	a.NoError(yaml.Unmarshal(byteValue, &tests))

	for i, t := range tests {
		result, err := getDesignSchema(t.Engine, t.Baseline, t.Target)
		a.NoError(err)
		if record {
			tests[i].Result = result
		} else {
			a.Equal(t.Result, result)
		}
	}

	if record {
		byteValue, err := yaml.Marshal(tests)
		a.NoError(err)
		err = os.WriteFile(filepath, byteValue, 0644)
		a.NoError(err)
	}
}

type checkTest struct {
	Engine   v1pb.Engine
	Metadata *v1pb.DatabaseMetadata
	Err      string
}

func TestCheckDatabaseMetadata(t *testing.T) {
	const (
		record = false
	)
	var (
		filepath = "testdata/check.yaml"
	)

	a := require.New(t)
	yamlFile, err := os.Open(filepath)
	a.NoError(err)

	tests := []checkTest{}
	byteValue, err := io.ReadAll(yamlFile)
	a.NoError(yamlFile.Close())
	a.NoError(err)
	a.NoError(yaml.Unmarshal(byteValue, &tests))

	for i, t := range tests {
		err := checkDatabaseMetadata(t.Engine, t.Metadata)
		if record {
			if err != nil {
				tests[i].Err = err.Error()
			} else {
				tests[i].Err = ""
			}
		} else {
			if t.Err == "" {
				a.NoError(err)
			} else {
				a.Equal(t.Err, err.Error())
			}
		}
	}

	if record {
		byteValue, err := yaml.Marshal(tests)
		a.NoError(err)
		err = os.WriteFile(filepath, byteValue, 0644)
		a.NoError(err)
	}
}
