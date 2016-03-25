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

func i64(i int) *int64 {
	j := int64(i)
	return &j
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func work(i int, pods kcl.PodInterface, namespace string, wg *sync.WaitGroup, t *time.Timer) {
	defer wg.Done()
	for {
		select {
		case <-t.C:
			log.Printf("Goroutine %d done, returning", i)
			return
		default:
		}
		podName := fmt.Sprintf("work-pod-%s", uuid.New())
		pod := &api.Pod{
			ObjectMeta: api.ObjectMeta{
				Name:      podName,
				Namespace: namespace,
			},
			Spec: api.PodSpec{
				Containers: []api.Container{
					api.Container{
						Name:            "alpine-echo",
						Image:           "alpine:3.3",
						Command:         []string{"echo", fmt.Sprintf(`"hello k8stress pod %s"`, podName)},
						ImagePullPolicy: api.PullIfNotPresent,
					},
				},
			},
		}
		newPod, err := pods.Create(pod)
		if err != nil {
			log.Printf("Error creating pod #%d %s (%s)", i, podName, err)
			return
		}
		log.Printf("New pod %s created: %+v", podName, *newPod)
		log.Printf("Deleting pod %s", podName)
		if err := pods.Delete(podName, &api.DeleteOptions{GracePeriodSeconds: i64(0)}); err != nil {
			log.Printf("Error deleting pod #%d %s (%s)", i, podName, err)
		}
		time.Sleep(1 * time.Second)
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
		go work(i, pods, conf.Namespace, &wg, timer)
	}
	log.Printf("Done creating %d goroutines", conf.NumGoroutines)
	wg.Wait()
	log.Printf("Done")
}
