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
	"github.com/stretchr/testify/require"
	"github.com/codeallergy/glue"
	"log"
	"reflect"
	"testing"
)

var reloadableBeanClass = reflect.TypeOf((*reloadableBean)(nil))

type reloadableBean struct {
	constructed int
	destroyed   int
}

func (t *reloadableBean) PostConstruct() error {
	t.constructed++
	return nil
}

func (t *reloadableBean) Destroy() error {
	t.destroyed++
	return nil
}

type topBean struct {
	ReloadableBean *reloadableBean `inject`
}

func TestBeanReload(t *testing.T) {

	reBean := &reloadableBean{}
	tBean := &topBean{}

	// initialization order
	ctx, err := glue.New(
		glue.Verbose{ Log: log.Default() },
		reBean,
		tBean,
	)
	require.NoError(t, err)

	require.Equal(t, 1, reBean.constructed)
	require.Equal(t, 0, reBean.destroyed)
	require.True(t, tBean.ReloadableBean == reBean)

	list := ctx.Bean(reloadableBeanClass, glue.DefaultLevel)
	require.Equal(t, 1, len(list))
	require.Equal(t, reBean, list[0].Object())

	err = list[0].Reload()
	require.NoError(t, err)

	require.Equal(t, 2, reBean.constructed)
	require.Equal(t, 1, reBean.destroyed)

	ctx.Close()

	require.Equal(t, 2, reBean.constructed)
	require.Equal(t, 2, reBean.destroyed)
	require.True(t, tBean.ReloadableBean == reBean)

}
