package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/garethjevans/build-action/pkg/logs"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const sleepTimeBetweenChecks = 3

var (
	v1alpha2Builds         = schema.GroupVersionResource{Group: "kpack.io", Version: "v1alpha2", Resource: "builds"}
	v1alpha2ClusterBuilder = schema.GroupVersionResource{Group: "kpack.io", Version: "v1alpha2", Resource: "clusterbuilders"}
)

func GetClusterBuilder(ctx context.Context, client dynamic.Interface, name string) (string, string, error) {
	clusterBuilder, err := client.Resource(v1alpha2ClusterBuilder).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}

	clusterBuilderName, _, err := unstructured.NestedString(clusterBuilder.Object, "status", "latestImage")
	if err != nil {
		return "", "", err
	}

	runImage, _, err := unstructured.NestedString(clusterBuilder.Object, "status", "stack", "runImage")
	if err != nil {
		return "", "", err
	}
	return clusterBuilderName, runImage, nil
}

func CreateBuild(ctx context.Context, client dynamic.Interface, namespace string, build *unstructured.Unstructured) (string, error) {
	created, err := client.Resource(v1alpha2Builds).Namespace(namespace).Create(ctx, build, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	return created.GetName(), nil
}

func GetBuild(ctx context.Context, client dynamic.Interface, namespace string, build string) (string, string, error) {
	got, err := client.Resource(v1alpha2Builds).Namespace(namespace).Get(ctx, build, metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}

	podName, _, err := unstructured.NestedString(got.Object, "status", "podName")
	if err != nil {
		return "", "", err
	}

	latestImage, _, err := unstructured.NestedString(got.Object, "status", "latestImage")
	if err != nil {
		return "", "", err
	}

	return podName, latestImage, nil
}

func main() {
	caCert := os.Getenv("CA_CERT")
	server := os.Getenv("SERVER")
	namespace := MustGetEnv("NAMESPACE")
	token := os.Getenv("TOKEN")

	gitRepo := fmt.Sprintf("%s/%s", MustGetEnv("GITHUB_SERVER_URL"), MustGetEnv("GITHUB_REPOSITORY"))
	gitSha := MustGetEnv("GITHUB_SHA")
	tag := MustGetEnv("TAG")

	decodedCaCert, err := base64.StdEncoding.DecodeString(caCert)
	if err != nil {
		panic(err)
	}

	var config *rest.Config

	if caCert == "" && server == "" && token == "" {
		// assume we are currently running inside the cluster we want to create the image resource in
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	} else {
		config = &rest.Config{
			TLSClientConfig: rest.TLSClientConfig{
				CAData: decodedCaCert,
			},
			Host:        server,
			BearerToken: token,
		}
	}

	ctx := context.Background()

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	clusterBuilder, runImage, err := GetClusterBuilder(ctx, dynamicClient, "default")
	if err != nil {
		panic(err)
	}

	build := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "kpack.io/v1alpha2",
			"kind":       "Build",
			"metadata": map[string]interface{}{
				"generateName": "gevans-petclinic-build-",
				"namespace":    namespace,
			},
			"spec": map[string]interface{}{
				"builder": map[string]interface{}{
					"image": clusterBuilder,
				},
				"runImage": map[string]interface{}{
					"image": runImage,
				},
				"serviceAccountName": "default",
				"source": map[string]interface{}{
					"git": map[string]interface{}{
						"url":      gitRepo,
						"revision": gitSha,
					},
				},
				"tags": []string{
					tag,
				},
			},
		},
	}

	name, err := CreateBuild(ctx, dynamicClient, namespace, build)
	if err != nil {
		panic(err)
	}

	for {
		var podName string
		podName, _, err = GetBuild(ctx, dynamicClient, namespace, name)
		if err != nil {
			panic(err)
		}

		if podName != "" {
			fmt.Printf("[DEBUG] Building... podName=%s\n", podName)
			err = GetPodLogs(ctx, client, namespace, podName)
			if err != nil {
				panic(err)
			}
			break
		}

		time.Sleep(sleepTimeBetweenChecks * time.Second)
	}

	fmt.Printf("[DEBUG] Checking build complete?")

	var latestImage string
	_, latestImage, err = GetBuild(ctx, dynamicClient, namespace, name)
	if err != nil {
		panic(err)
	}

	if latestImage != "" {
		fmt.Printf("::set-output name=name::%s\n", latestImage)
	}
}

func GetPodLogs(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, podName string) error {
	st := logs.SternTailer{}
	return st.Tail(ctx, clientSet, namespace, podName)
}

func MustGetEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("Environment Var %s must be set", name)
	}
	return val
}
