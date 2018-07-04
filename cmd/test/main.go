package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	midonet "github.com/midonet/midonet-kubernetes/pkg/client/clientset/versioned"
)

func main() {
	flag.Parse()

	kubeconfig := ".kube/config"
	k8sconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.WithError(err).Fatal("BuildConfigFromFlags")
	}

	mnClientset, err := midonet.NewForConfig(k8sconfig)
	if err != nil {
		log.WithError(err).Fatal("midonet.NewForConfig")
	}

	opts := metav1.ListOptions{}
	list, err := mnClientset.MidonetV1().Translations("kube-system").List(opts)
	if err != nil {
		log.WithError(err).Fatal("List")
	}
	log.WithField("list", list).Info("Get a list")
	log.WithField("items", list.Items).Info("Items")
}
