// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"fmt"
	"runtime"
)

var (
	VERSION     string
	BUILD_TIME  string
	COMMIT_HASH string
	COMMIT_TIME string
)

func GetVersion() string {
	if VERSION != "" {
		return fmt.Sprintf(
			"\n\tfree5GC version: %s"+
				"\n\tbuild time:      %s"+
				"\n\tcommit hash:     %s"+
				"\n\tcommit time:     %s"+
				"\n\tgo version:      %s %s/%s",
			VERSION,
			BUILD_TIME,
			COMMIT_HASH,
			COMMIT_TIME,
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH,
		)
	} else {
		return fmt.Sprintf(
			"\n\tNot specify ldflags (which link version) during go build\n\tgo version: %s %s/%s",
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH,
		)
	}
}
