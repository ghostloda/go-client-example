package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yaml2 "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"os"
)

func ToCreatedDeployment() {

	clientset := newClient()
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	// List Deployments
	prompt()
	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	// Delete Deployment
	prompt()
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete("demo-deployment", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")
}

// yaml文件转json,然后反序列化 type deployment struct

func yamlCreateDeployment() {
	var (
		clientset  *kubernetes.Clientset
		deployYaml []byte
		deployJson []byte
		deployment = appsv1.Deployment{}
		replicas   int32
		err        error
	)

	// 初始化k8s客户端
	clientset = newClient()

	// 读取YAML
	if deployYaml, err = ioutil.ReadFile("./nginx.yaml"); err != nil {
		goto FAIL
	}

	// YAML转JSON
	if deployJson, err = yaml2.ToJSON(deployYaml); err != nil {
		goto FAIL
	}

	// JSON转struct
	if err = json.Unmarshal(deployJson, &deployment); err != nil {
		goto FAIL
	}

	// 修改replicas数量为1
	replicas = 1
	deployment.Spec.Replicas = &replicas

	// 查询k8s是否有该deployment
	if _, err = clientset.AppsV1().Deployments("default").Get(deployment.Name, metav1.GetOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			goto FAIL
		}
		// 不存在则创建
		if _, err = clientset.AppsV1().Deployments("default").Create(&deployment); err != nil {
			goto FAIL
		}
	} else { // 已存在则更新
		if _, err = clientset.AppsV1().Deployments("default").Update(&deployment); err != nil {
			goto FAIL
		}
	}

	fmt.Println("apply成功!")
	return

FAIL:
	fmt.Println(err)
	return
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func int32Ptr(i int32) *int32 { return &i }
