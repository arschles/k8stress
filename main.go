package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

func main() {
	conf := new(Config)
	if err := envconfig.Process("k8stress", conf); err != nil {
		log.Fatalf("Error processing config (%s)", err)
	}
	host, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error getting hostname (%s)", err)
	}
	client, err := kcl.NewInCluster()
	if err != nil {
		log.Fatalf("Error creating new Kubernetes client (%s)", err)
	}
	log.Printf("Running on host %s with config %s", host, conf)
	pods := client.Pods(conf.Namespace)
	var wg sync.WaitGroup
	for i := 0; i < conf.NumGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			pod := &api.Pod{
				ObjectMeta: api.ObjectMeta{
					Name:      fmt.Sprintf("k8stress-pod-%s-%d", host, i),
					Namespace: conf.Namespace,
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						api.Container{
							Name:            "alpine-echo",
							Image:           "alpine:3.3",
							Command:         []string{"echo", fmt.Sprintf(`"hello k8stress pod %d"`, i)},
							ImagePullPolicy: api.PullIfNotPresent,
						},
					},
				},
			}
			newPod, err := pods.Create(pod)
			if err != nil {
				log.Printf("Error creating pod %d (%s)", i, err)
				return
			}
			log.Printf("New pod #%d created: %+v", i, *newPod)
		}(i)
	}
	wg.Wait()
	log.Printf("Done creating %d pods", conf.NumGoroutines)
}
