/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package glue

import "log"

/**
Verbose logs if not nil
*/
var verbose *log.Logger

/**
Use this function operate verbose and logging level during context creation.
*/

func Verbose(log *log.Logger) (prev *log.Logger) {
	prev, verbose = verbose, log
	return
}

