package main

import (
	"context"
	"os"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func login() (*kubernetes.Clientset, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func doWork() {
	time.Sleep(2 * time.Second)
	logger.Info("work done!")
}

func createJob(cs *kubernetes.Clientset) error {
	jobs := clientset.BatchV1().Jobs(os.Getenv("POD_NS"))
	jobName := "do-work"
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: os.Getenv("POD_NS"),
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    jobName,
							Image:   "localhost:5000/api:latest",
							Command: []string{"./app", "--job"},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
		},
	}

	_, err := jobs.Create(context.Background(), jobSpec, metav1.CreateOptions{})

	return err
}
