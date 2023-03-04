/*
 * Copyright (c) 2022-2023 Zander Schwid & Co. LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */

package glue_test

import (
	"errors"
	"github.com/codeallergy/glue"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"strings"
	"testing"
)

type fileSystemStub struct {
}

func (t fileSystemStub) Open(name string) (http.File, error) {
	return nil, errors.New(name)
}

func TestResourceMerge(t *testing.T) {

	ctx, err := glue.New(
		glue.Verbose{ Log: log.Default() },
		glue.ResourceSource{
			Name: "resources",
			AssetNames: []string{ "a.txt", "b/c.txt" },
			AssetFiles: fileSystemStub{},
		},
		&glue.ResourceSource{
			Name: "resources",
			AssetNames: []string{ "d.txt", "f/g.txt" },
			AssetFiles: fileSystemStub{},
		},
	)

	require.NoError(t, err)
	defer ctx.Close()

	validNames := []string {
		"a.txt",
		"b/c.txt",
		"d.txt",
		"f/g.txt",
	}

	for _, validName := range validNames {
		res, ok := ctx.Resource("resources:" + validName)
		require.True(t, ok)
		_, err = res.Open()
		require.Equal(t, validName, err.Error())
	}

	_, ok := ctx.Resource("assets:a.txt")
	require.False(t, ok)

	_, ok = ctx.Resource("a.txt")
	require.False(t, ok)

}

func TestResourceMergeConflict(t *testing.T) {

	ctx, err := glue.New(
		glue.Verbose{ Log: log.Default() },
		glue.ResourceSource{
			Name: "resources",
			AssetNames: []string{ "a.txt", "b/c.txt" },
			AssetFiles: fileSystemStub{},
		},
		&glue.ResourceSource{
			Name: "resources",
			AssetNames: []string{ "a.txt" },
			AssetFiles: fileSystemStub{},
		},
	)

	require.Error(t, err)
	require.Nil(t, ctx)
	println(err.Error())
	require.True(t, strings.Contains(err.Error(), "already exist"))

}

func TestResourceParent(t *testing.T) {

	parent, err := glue.New(
		glue.Verbose{ Log: log.Default() },
		glue.ResourceSource{
			Name: "resources",
			AssetNames: []string{ "a.txt", "b/c.txt" },
			AssetFiles: fileSystemStub{},
		},
	)

	require.NoError(t, err)
	defer parent.Close()

	child, err := parent.Extend(
		glue.ResourceSource{
			Name: "resources",
			AssetNames: []string{ "a.txt" },
			AssetFiles: fileSystemStub{},
		},
		&glue.ResourceSource{
			Name: "resources",
			AssetNames: []string{ "d.txt", "f/g.txt" },
			AssetFiles: fileSystemStub{},
		},
	)
	require.NoError(t, err)
	defer child.Close()

	validNames := []string {
		"a.txt",
		"b/c.txt",
		"d.txt",
		"f/g.txt",
	}

	for _, validName := range validNames {
		res, ok := child.Resource("resources:" + validName)
		require.True(t, ok)
		_, err = res.Open()
		require.Equal(t, validName, err.Error())
	}

}