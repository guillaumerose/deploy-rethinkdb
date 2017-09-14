package main

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	configFile, err := homedir.Expand("~/.kube/config")
	if err != nil {
		panic(err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", configFile)
	if err != nil {
		panic(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	_, err = client.CoreV1().Secrets(v1.NamespaceDefault).Create(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "dtr-rethinkdb",
		},
		StringData: map[string]string{
			"rethinkdb-password": base64.StdEncoding.EncodeToString([]byte("cmV0aGlua2Ri")),
		},
	})
	if err != nil {
		panic(err)
	}

	if err = createService(client, "service1-rethinkdb.yaml"); err != nil {
		panic(err)
	}
	if err = createService(client, "service2-rethinkdb.yaml"); err != nil {
		panic(err)
	}
	if err = createService(client, "service3-rethinkdb.yaml"); err != nil {
		panic(err)
	}

	if err = createDeployment(client, "dep1-rethinkdb.yaml"); err != nil {
		panic(err)
	}
	if err = createStatefulSet(client, "dep2-rethinkdb.yaml"); err != nil {
		panic(err)
	}
}

func createStatefulSet(client *kubernetes.Clientset, file string) error {
	bin, _ := ioutil.ReadFile(file)
	var dep v1beta1.StatefulSet
	err := yaml.Unmarshal(bin, &dep)
	if err != nil {
		return err
	}
	_, err = client.AppsV1beta1().StatefulSets(v1.NamespaceDefault).Create(&dep)
	return err
}

func createDeployment(client *kubernetes.Clientset, file string) error {
	bin, _ := ioutil.ReadFile(file)
	var dep v1beta1.Deployment
	err := yaml.Unmarshal(bin, &dep)
	if err != nil {
		return err
	}
	_, err = client.AppsV1beta1().Deployments(v1.NamespaceDefault).Create(&dep)
	return err
}

func createService(client *kubernetes.Clientset, file string) error {
	bin, _ := ioutil.ReadFile(file)
	var service v1.Service
	err := yaml.Unmarshal(bin, &service)
	if err != nil {
		return err
	}
	_, err = client.CoreV1().Services(v1.NamespaceDefault).Create(&service)
	return err
}
