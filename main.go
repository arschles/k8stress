package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

func main() {
	conf := new(Config)
	if err := envconfig.Process("k8stress", conf); err != nil {
		log.Fatalf("Error processing config (%s)", err)
	}
	client, err := kcl.NewInCluster()
	if err != nil {
		log.Fatalf("Error creating new Kubernetes client (%s)", err)
	}
	pods := client.Pods(conf.Namespace)
	for i := 0; i < conf.NumGoroutines; i++ {
		go func(i int) {
			pod := &api.Pod{}
			newPod, err := pods.Create(pod)
			if err != nil {
				log.Printf("Error creating pod %d (%s)", i, err)
				return
			}
			log.Printf("New pod #%d created:\n%+v", i, *newPod)
		}(i)
	}
}
