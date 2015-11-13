// Copyright 2015 clair authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package os

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/coreos/clair/worker/detectors"
)

var (
	osReleaseOSRegexp      = regexp.MustCompile(`^ID=(.*)`)
	osReleaseVersionRegexp = regexp.MustCompile(`^VERSION_ID=(.*)`)
)

// OsReleaseOSDetector implements OSDetector and detects the OS from the
// /etc/os-release and usr/lib/os-release files.
type OsReleaseOSDetector struct{}

func init() {
	detectors.RegisterOSDetector("os-release", &OsReleaseOSDetector{})
}

// Detect tries to detect OS/Version using "/etc/os-release" and "/usr/lib/os-release"
// Typically for Debian / Ubuntu
// /etc/debian_version can't be used, it does not make any difference between testing and unstable, it returns stretch/sid
func (detector *OsReleaseOSDetector) Detect(data map[string][]byte) (OS, version string) {
	for _, filePath := range detector.GetRequiredFiles() {
		f, hasFile := data[filePath]
		if !hasFile {
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(string(f)))
		for scanner.Scan() {
			line := scanner.Text()

			r := osReleaseOSRegexp.FindStringSubmatch(line)
			if len(r) == 2 {
				OS = strings.Replace(strings.ToLower(r[1]), "\"", "", -1)
			}

			r = osReleaseVersionRegexp.FindStringSubmatch(line)
			if len(r) == 2 {
				version = strings.Replace(strings.ToLower(r[1]), "\"", "", -1)
			}
		}
	}

	return
}

// GetRequiredFiles returns the list of files that are required for Detect()
func (detector *OsReleaseOSDetector) GetRequiredFiles() []string {
	return []string{"etc/os-release", "usr/lib/os-release"}
}