// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compliance

import (
	"bytes"
	"testing"
)

func TestResolveBottomUpConditions(t *testing.T) {
	tests := []struct {
		name                string
		roots               []string
		edges               []annotated
		expectedResolutions []res
	}{
		{
			name:  "firstparty",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name:  "firstpartytool",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"toolchain"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name:  "firstpartydeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name:  "firstpartywide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name:  "firstpartydynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name:  "firstpartydynamicdeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name:  "firstpartydynamicwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name:  "restricted",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "restrictedtool",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"toolchain"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "restricteddeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "restrictedwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "restricteddynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "restricteddynamicdeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "restricteddynamicwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "weakrestricted",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "weakrestrictedtool",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"toolchain"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "weakrestricteddeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "weakrestrictedwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "weakrestricteddynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "weakrestricteddynamicdeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "weakrestricteddynamicwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "classpath",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name:  "classpathdependent",
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name:  "classpathdynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name:  "classpathdependentdynamic",
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			lg, err := toGraph(stderr, tt.roots, tt.edges)
			if err != nil {
				t.Errorf("unexpected test data error: got %w, want no error", err)
				return
			}
			expectedRs := toResolutionSet(lg, tt.expectedResolutions)
			actualRs := ResolveBottomUpConditions(lg)
			checkSame(actualRs, expectedRs, t)
		})
	}
}

func TestResolveTopDownConditions(t *testing.T) {
	tests := []struct {
		name                string
		roots               []string
		edges               []annotated
		expectedResolutions []res
	}{
		{
			name:  "firstparty",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name:  "firstpartydynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name:  "restricted",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name:  "restrictedtool",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplBin.meta_lic", []string{"toolchain"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"gplBin.meta_lic", "gplBin.meta_lic", "gplBin.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name:  "restricteddeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "mitBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "mplLib.meta_lic", []string{"static"}},
				{"mitBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "mplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "mplLib.meta_lic", "mplLib.meta_lic", "reciprocal"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mplLib.meta_lic", "mplLib.meta_lic", "reciprocal"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitBin.meta_lic", "mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"mitBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"mplLib.meta_lic", "mplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mplLib.meta_lic", "mplLib.meta_lic", "mplLib.meta_lic", "reciprocal"},
			},
		},
		{
			name:  "restrictedwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "restricteddynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name:  "restricteddynamicdeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "mitBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
				{"apacheBin.meta_lic", "mplLib.meta_lic", []string{"dynamic"}},
				{"mitBin.meta_lic", "mitLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "mplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitBin.meta_lic", "mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"mplLib.meta_lic", "mplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mplLib.meta_lic", "mplLib.meta_lic", "mplLib.meta_lic", "reciprocal"},
			},
		},
		{
			name:  "restricteddynamicwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "weakrestricted",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name:  "weakrestrictedtool",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplBin.meta_lic", []string{"toolchain"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"lgplBin.meta_lic", "lgplBin.meta_lic", "lgplBin.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name:  "weakrestricteddeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "mitLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name:  "weakrestrictedwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "weakrestricteddynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name:  "weakrestricteddynamicdeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "weakrestricteddynamicwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name:  "classpath",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name:  "classpathdependent",
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
				{"dependentModule.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"dependentModule.meta_lic", "mitLib.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name:  "classpathdynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name:  "classpathdependentdynamic",
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
				{"dependentModule.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"dependentModule.meta_lic", "mitLib.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			lg, err := toGraph(stderr, tt.roots, tt.edges)
			if err != nil {
				t.Errorf("unexpected test data error: got %w, want no error", err)
				return
			}
			expectedRs := toResolutionSet(lg, tt.expectedResolutions)
			actualRs := ResolveTopDownConditions(lg)
			checkSame(actualRs, expectedRs, t)
		})
	}
}
