package main

import (
	"context"
	"flag"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	return os.Getenv("USERPROFILE") // windows
}

func connectToK8s() *kubernetes.Clientset {
	var kubeconfig *string
	home := homeDir()
	fmt.Println("home : " + home)

	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	fmt.Println("config path : " + *kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Panicln("failed to create k8s config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln("failed to create k8s clientset")
	}

	return clientset
}

func launchK8sJob(clientset *kubernetes.Clientset, jobName *string, image *string, cmd *string) {
	jobs := clientset.BatchV1().Jobs("default")
	var backOffLimit int32 = 0

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *jobName,
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    *jobName,
							Image:   *image,
							Command: strings.Split(*cmd, " "),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln("Failed to create K8s job.")
	}

	//print job details
	log.Println("Created K8s job successfully")
}
:q
func main() {
	fmt.Println("v1.1.2")
	jobName := flag.String("jobname", "test-job", "The name of the job")
	containerImage := flag.String("image", "ubuntu:latest", "Name of the container image")
	entryCommand := flag.String("command", "ls", "The command to run inside the container")

	flag.Parse()

	fmt.Printf("Args : %s %s %s \n", *jobName, *containerImage, *entryCommand)

	clientset := connectToK8s()
	launchK8sJob(clientset, jobName, containerImage, entryCommand)
}
