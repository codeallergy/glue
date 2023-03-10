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
	"testing"
)

/**
Cycle dependency test of plain beans
*/

type aPlainBean struct {
	BBean *bPlainBean `inject`
}

type bPlainBean struct {
	CBean *cPlainBean `inject`
}

type cPlainBean struct {
	ABean *aPlainBean `inject:"lazy"`
}

func TestPlainBeanCycle(t *testing.T) {

	ctx, err := glue.New(
		glue.Verbose{ Log: log.Default() },
		&aPlainBean{},
		&bPlainBean{},
		&cPlainBean{},
	)
	require.NoError(t, err)
	defer ctx.Close()

}

type selfDepBean struct {
	Self *selfDepBean `inject`
}

func TestSelfDepCycle(t *testing.T) {

	self := &selfDepBean{}

	ctx, err := glue.New(
		glue.Verbose{ Log: log.Default() },
		self,
	)
	require.NoError(t, err)
	defer ctx.Close()

	require.True(t, self == self.Self)

}