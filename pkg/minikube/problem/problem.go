/*
Copyright 2019 The Kubernetes Authors All rights reserved.

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

// Package problem helps deliver actionable feedback to a user based on an error message.
package problem

import (
	"regexp"

	"k8s.io/minikube/pkg/minikube/console"
)

const issueBase = "https://github.com/kubernetes/minikube/issue"

// Problem represents a known problem in minikube.
type Problem struct {
	ID     string
	Err    error
	Advice string
	URL    string
	Issues []int
}

// match maps a regular expression to problem metadata.
type match struct {
	Regexp *regexp.Regexp
	Advice string
	URL    string
	Issues []int
}

// Display problem metadata to the console
func (p *Problem) Display() {
	console.ErrStyle("failure", "Error:         [%s] %v", p.ID, p.Err)
	console.ErrStyle("tip", "Advice:        %s", p.Advice)
	if p.URL != "" {
		console.ErrStyle("documentation", "Documentation: %s", p.URL)
	}
	if len(p.Issues) == 0 {
		return
	}
	console.ErrStyle("issues", "Related issues:")
	issues := p.Issues
	if len(issues) > 3 {
		issues = issues[0:3]
	}
	for _, i := range issues {
		console.ErrStyle("issue", "%s/%d", issueBase, i)
	}
}

// FromError returns a known problem from an error.
func FromError(err error) *Problem {
	maps := []map[string]match{
		osProblems,
		vmProblems,
		netProblems,
		deployProblems,
	}
	for _, m := range maps {
		for k, v := range m {
			if v.Regexp.MatchString(err.Error()) {
				return &Problem{
					Err:    err,
					Advice: v.Advice,
					URL:    v.URL,
					ID:     k,
					Issues: v.Issues,
				}
			}
		}
	}
	return nil
}
