/*
Copyright 2021 VMware, Inc.

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

package logs

import (
	"context"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"os"
	"regexp"
	"text/template"
	"time"

	"github.com/fatih/color"
	"github.com/stern/stern/stern"
	"k8s.io/apimachinery/pkg/fields"
)

type SternTailer struct {
}

func (s *SternTailer) Tail(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, podName string) error {
	t := "{{color .PodColor \"[\"}}{{color .PodColor .ContainerName}}{{color .PodColor \"]\"}} {{.Message}}\n"

	functions := map[string]interface{}{
		"color": func(color color.Color, text string) string {
			return color.SprintFunc()(text)
		},
	}
	parsedTemplate, err := template.New("log").Funcs(functions).Parse(t)
	if err != nil {
		panic(err)
	}

	configStern := stern.Config{
		Namespaces:     []string{namespace},
		Location:       time.Local,
		LabelSelector:  labels.Everything(),
		ContainerQuery: regexp.MustCompile(".*"),
		ContainerStates: []stern.ContainerState{
			stern.RUNNING,
			stern.TERMINATED,
		},
		InitContainers: true,
		Since:          10 * time.Second,
		PodQuery:       regexp.MustCompile(podName),
		FieldSelector:  fields.Everything(),
		Template:       parsedTemplate,
		Out:            os.Stdout,
		ErrOut:         os.Stderr,
	}

	return Run(ctx, clientSet, &configStern)
}
