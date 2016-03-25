package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pborman/uuid"
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func work(i int, pods kcl.PodInterface, podName, namespace string, wg *sync.WaitGroup, t *time.Timer) {
	defer wg.Done()
	for {
		select {
		case <-t.C:
			return
		default:
		}
		pod := &api.Pod{
			ObjectMeta: api.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", podName, uuid.New()),
				Namespace: namespace,
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
	}
}

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
	timer := time.NewTimer(time.Duration(conf.TimeSec) * time.Second)
	defer timer.Stop()
	for i := 0; i < conf.NumGoroutines; i++ {
		wg.Add(1)
		podName := fmt.Sprintf("k8stress-pod-%s-%d", host, i)
		go work(i, pods, podName, conf.Namespace, &wg, timer)
	}
	log.Printf("Done creating %d goroutines", conf.NumGoroutines)
	wg.Wait()
	log.Printf("Done")
}
