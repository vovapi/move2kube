/*
Copyright IBM Corporation 2020

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package output_test

import (
	"reflect"
	"testing"

	"github.com/konveyor/move2kube/types/output"
)

func TestMerge(t *testing.T) {
	t.Run("merge 2 empty helm values", func(t *testing.T) {
		h1 := output.HelmValues{}
		h2 := output.HelmValues{}
		want := output.HelmValues{}
		if h1.Merge(h2); !reflect.DeepEqual(h1, want) {
			t.Fatal("The value should not have changed after merge. Expected:", want, "Actual:", h1)
		}
	})

	t.Run("merge filled helm value into filled helm value", func(t *testing.T) {
		h1 := output.HelmValues{}
		h1.RegistryNamespace = "namespace1"
		h1.RegistryURL = "url1"
		h1.StorageClass = "storagecls1"
		h2 := output.HelmValues{}
		h2.RegistryNamespace = "namespace2"
		h2.RegistryURL = "url2"
		h2.StorageClass = "storagecls2"
		want := output.HelmValues{}
		want.RegistryNamespace = "namespace2"
		want.RegistryURL = "url2"
		want.StorageClass = "storagecls2"
		if h1.Merge(h2); !reflect.DeepEqual(h1, want) {
			t.Fatal("Failed to merge the helm values properly. Expected:", want, "Actual:", h1)
		}
	})

	t.Run("merge global and service variables into filled helm value", func(t *testing.T) {
		makeH := func() output.HelmValues {
			h := output.HelmValues{}
			h.GlobalVariables = make(map[string]string)
			return h
		}
		key1 := "key1"
		val1 := "val1"
		val2 := "val2"

		h1 := makeH()
		h1.GlobalVariables[key1] = val1

		h2 := makeH()
		h2.GlobalVariables[key1] = val2

		want := makeH()
		want.GlobalVariables[key1] = val2

		if h1.Merge(h2); !reflect.DeepEqual(h1, want) {
			t.Fatal("Failed to merge the helm values properly. Expected:", want, "Actual:", h1)
		}
	})

	t.Run("merge ImageTagTree properly into filled helm value", func(t *testing.T) {
		makeH := func() output.HelmValues {
			h := output.HelmValues{}
			h.Services = make(map[string]output.Service)
			return h
		}
		key1 := "key1"
		key2 := "key2"
		con1 := "name1"
		val1 := output.Container{"tag1"}
		val2 := output.Container{"tag2"}

		h1 := makeH()
		h1.Services[key1] = output.Service{map[string]output.Container{con1: val1}}

		h2 := makeH()
		h2.Services[key1] = output.Service{map[string]output.Container{con1: val2}}
		h2.Services[key2] = output.Service{map[string]output.Container{con1: val1}}

		want := makeH()
		want.Services[key1] = output.Service{map[string]output.Container{con1: val2}}
		want.Services[key2] = output.Service{map[string]output.Container{con1: val1}}

		if h1.Merge(h2); !reflect.DeepEqual(h1, want) {
			t.Fatal("Failed to merge the helm values properly. Expected:", want, "Actual:", h1)
		}
	})
}
