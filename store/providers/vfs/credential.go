/*
Copyright The Pharmer Authors.

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
package vfs

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	cloudapi "pharmer.dev/cloud/pkg/apis/cloud/v1"
	api "pharmer.dev/pharmer/apis/v1alpha1"
	"pharmer.dev/pharmer/store"

	"github.com/pkg/errors"
	"gomodules.xyz/stow"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type credentialFileStore struct {
	container stow.Container
	prefix    string
}

var _ store.CredentialStore = &credentialFileStore{}

func (s *credentialFileStore) resourceHome() string {
	return filepath.Join(s.prefix, "credentials")
}

func (s *credentialFileStore) resourceID(name string) string {
	return filepath.Join(s.resourceHome(), name+".json")
}

func (s *credentialFileStore) List(opts metav1.ListOptions) ([]*cloudapi.Credential, error) {
	result := make([]*cloudapi.Credential, 0)
	cursor := stow.CursorStart
	for {
		page, err := s.container.Browse(s.resourceHome()+"/", string(os.PathSeparator), cursor, pageSize)
		if err != nil {
			return nil, errors.Errorf("failed to list credentials. Reason: %v", err)
		}
		for _, item := range page.Items {
			r, err := item.Open()
			if err != nil {
				return nil, errors.Errorf("failed to list credentials. Reason: %v", err)
			}
			var obj cloudapi.Credential
			err = json.NewDecoder(r).Decode(&obj)
			if err != nil {
				return nil, errors.Errorf("failed to list credentials. Reason: %v", err)
			}
			result = append(result, &obj)
			r.Close()
		}
		cursor = page.Cursor
		if stow.IsCursorEnd(cursor) {
			break
		}
	}
	return result, nil
}

func (s *credentialFileStore) Get(name string) (*cloudapi.Credential, error) {
	if name == "" {
		return nil, errors.New("missing credential name")
	}
	item, err := s.container.Item(s.resourceID(name))
	if err != nil {
		return nil, errors.Errorf("credential `%s` does not exist. Reason: %v", name, err)
	}

	r, err := item.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var existing cloudapi.Credential
	err = json.NewDecoder(r).Decode(&existing)
	if err != nil {
		return nil, err
	}
	return &existing, nil
}

func (s *credentialFileStore) Create(obj *cloudapi.Credential) (*cloudapi.Credential, error) {
	if obj == nil {
		return nil, errors.New("missing credential")
	} else if obj.Name == "" {
		return nil, errors.New("missing credential name")
	}
	err := api.AssignTypeKind(obj)
	if err != nil {
		return nil, err
	}

	id := s.resourceID(obj.Name)
	_, err = s.container.Item(id)
	if err == nil {
		return nil, errors.Errorf("credential `%s` already exists", obj.Name)
	}

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, err
	}
	_, err = s.container.Put(id, bytes.NewBuffer(data), int64(len(data)), nil)
	return obj, err
}

func (s *credentialFileStore) Update(obj *cloudapi.Credential) (*cloudapi.Credential, error) {
	if obj == nil {
		return nil, errors.New("missing credential")
	} else if obj.Name == "" {
		return nil, errors.New("missing credential name")
	}
	err := api.AssignTypeKind(obj)
	if err != nil {
		return nil, err
	}

	id := s.resourceID(obj.Name)

	_, err = s.container.Item(id)
	if err != nil {
		return nil, errors.Errorf("credential `%s` does not exist. Reason: %v", obj.Name, err)
	}

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, err
	}
	_, err = s.container.Put(id, bytes.NewBuffer(data), int64(len(data)), nil)
	return obj, err
}

func (s *credentialFileStore) Delete(name string) error {
	if name == "" {
		return errors.New("missing credential name")
	}

	resourceID := s.resourceID(name)

	item, err := s.container.Item(resourceID)
	if err != nil {
		return errors.Errorf("failed to get item %s. Reason: %v", name, err)
	}

	return s.container.RemoveItem(item.ID())
}
