/*
Copyright 2016 The Kubernetes Authors.

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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/buger/jsonparser"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	client, _ := kubernetes.NewForConfig(config)
	clusterOperatorNameList := getAllClusterOperators(client)
	for _, name := range clusterOperatorNameList {
		fmt.Println(name)
	}
}

func getAllClusterOperators(client *kubernetes.Clientset) []string {

	var clusterOperatorNameList []string

	rawResponse, _ := client.RESTClient().
		Get().
		AbsPath("/apis/config.openshift.io/v1").
		Resource("clusteroperators").
		DoRaw(context.Background())

	jsonparser.ArrayEach(rawResponse, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		clusterOperatorName, _, _, _ := jsonparser.Get(value, "metadata", "name")
		clusterOperatorNameList = append(clusterOperatorNameList, string(clusterOperatorName))
	}, "items")
	return clusterOperatorNameList
}
