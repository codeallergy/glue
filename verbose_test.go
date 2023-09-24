/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package glue_test

import (
	"github.com/codeallergy/glue"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func init() {
	glue.Verbose(log.Default())
}

func TestVerbose(t *testing.T) {
	prev := glue.Verbose(log.Default())
	require.NotNil(t, prev)
}

