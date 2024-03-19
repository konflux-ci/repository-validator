package v1alpha1

import (
	"errors"
	"reflect"
	"testing"

	pacv1alpha1 "github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"

	corev1 "k8s.io/api/core/v1"
)

func TestLoadUrlPrefixAllowListFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		args       args
		fileReader FileReader
		want       []string
		wantErr    bool
	}{
		{
			name: "Load empty file",
			args: args{path: "file"},
			fileReader: func(name string) ([]byte, error) {
				return []byte("[]"), nil
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "Load non empty file",
			args: args{path: "file"},
			fileReader: func(name string) ([]byte, error) {
				return []byte(`
					[
						"a",
						"b",
						"b"
					]
				`), nil
			},
			want: []string{
				"a",
				"b",
				"b",
			},
			wantErr: false,
		},
		{
			name: "The given path is an empty string",
			args: args{path: ""},
			fileReader: func(name string) ([]byte, error) {
				return nil, nil
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "Load file with broken json",
			args: args{path: "file"},
			fileReader: func(name string) ([]byte, error) {
				return []byte("abc"), nil
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Load file random error",
			args: args{path: "file"},
			fileReader: func(name string) ([]byte, error) {
				return nil, errors.New("Random Error")
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadUrlPrefixAllowListFromFile(tt.args.path, tt.fileReader)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadUrlPrefixAllowListFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadUrlPrefixAllowListFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_castToRepository(t *testing.T) {
	type args struct {
		obj runtime.Object
	}
	tests := []struct {
		name   string
		args   args
		want   *pacv1alpha1.Repository
		wantOk bool
	}{
		{
			name: "Got non repository object",
			args: args{
				obj: &corev1.Pod{},
			},
			want:   nil,
			wantOk: false,
		},
		{
			name: "Got repository object",
			args: args{
				obj: &pacv1alpha1.Repository{},
			},
			want:   &pacv1alpha1.Repository{},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, wantOk := castToRepository(tt.args.obj)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("castToRepository() got = %v, want %v", got, tt.want)
			}
			if wantOk != tt.wantOk {
				t.Errorf("castToRepository() got1 = %v, want %v", wantOk, tt.wantOk)
			}
		})
	}
}
