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
	/*
		rawResponse, _ := client.RESTClient().
			Get().
			AbsPath("/apis/config.openshift.io/v1").
			Resource("clusteroperators").
			DoRaw(context.Background())
	*/

	/*
		rawResponse, _ := client.RESTClient().
			Get().
			AbsPath("/apis/metrics.k8s.io/v1beta1").
			Resource("pods").Namespace("observatorium").
			DoRaw(context.Background())
	*/

	rawResponse, _ := client.RESTClient().
		Get().
		AbsPath("/apis/monitoring.coreos.com/v1/").
		Resource("prometheuses").Namespace("openshift-monitoring").
		DoRaw(context.Background())
	fmt.Println(jsonparser.GetUnsafeString(rawResponse))
	/*
		jsonparser.ArrayEach(rawResponse, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			clusterOperatorName, _, _, _ := jsonparser.Get(value, "metadata", "name")
			clusterOperatorNameList = append(clusterOperatorNameList, string(clusterOperatorName))
		}, "items")
	*/
	return clusterOperatorNameList
}
