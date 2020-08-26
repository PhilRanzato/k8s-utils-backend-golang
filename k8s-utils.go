package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// Getter
func getConfig() *rest.Config {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	// log.Println("Using kubeconfig file: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func getClientSet() *kubernetes.Clientset {
	var config = getConfig()

	// Create a rest client not targeting specific API version
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientset
}

func getPods(clientset *kubernetes.Clientset) []v1.Pod {
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get pods:", err)
	}
	return pods.Items
}

func getDeployments(clientset *kubernetes.Clientset) []appv1.Deployment {
	deploys, err := clientset.AppsV1().Deployments("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get pods:", err)
	}
	return deploys.Items
}

func getDaemonSets(clientset *kubernetes.Clientset) []appv1.DaemonSet {
	ds, err := clientset.AppsV1().DaemonSets("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get pods:", err)
	}
	return ds.Items
}

func getReplicaSets(clientset *kubernetes.Clientset) []appv1.ReplicaSet {
	rs, err := clientset.AppsV1().ReplicaSets("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get pods:", err)
	}
	return rs.Items
}

func getStatefulSets(clientset *kubernetes.Clientset) []appv1.StatefulSet {
	sts, err := clientset.AppsV1().StatefulSets("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get pods:", err)
	}
	return sts.Items
}

func getServices(clientset *kubernetes.Clientset) []v1.Service {
	svcs, err := clientset.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}
	return svcs.Items
}

func getSecrets(clientset *kubernetes.Clientset) []v1.Secret {
	scr, err := clientset.CoreV1().Secrets("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}
	return scr.Items
}

func getConfigMaps(clientset *kubernetes.Clientset) []v1.ConfigMap {
	cm, err := clientset.CoreV1().ConfigMaps("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}
	return cm.Items
}

func getServiceAccounts(clientset *kubernetes.Clientset) []v1.ServiceAccount {
	sa, err := clientset.CoreV1().ServiceAccounts("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}
	return sa.Items
}

func getEvents(clientset *kubernetes.Clientset) []v1.Event {
	ev, err := clientset.CoreV1().Events("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}
	return ev.Items
}

func getEndpoints(clientset *kubernetes.Clientset) []v1.Endpoints {
	ep, err := clientset.CoreV1().Endpoints("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}
	return ep.Items
}

func getPersistentVolumes(clientset *kubernetes.Clientset) []v1.PersistentVolume {
	pvs, err := clientset.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}
	return pvs.Items
}

func getPersistentVolumeClaims(clientset *kubernetes.Clientset) []v1.PersistentVolumeClaim {
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}
	return pvcs.Items
}

func getRoles(clientset *kubernetes.Clientset) []rbacv1.Role {
	roles, err := clientset.RbacV1().Roles("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}
	return roles.Items
}

func getRoleBindings(clientset *kubernetes.Clientset) []rbacv1.RoleBinding {
	rb, err := clientset.RbacV1().RoleBindings("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}
	return rb.Items
}

func getClusterRoles(clientset *kubernetes.Clientset) []rbacv1.ClusterRole {
	cr, err := clientset.RbacV1().ClusterRoles().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}
	return cr.Items
}

func getClusterRoleBindings(clientset *kubernetes.Clientset) []rbacv1.ClusterRoleBinding {
	crb, err := clientset.RbacV1().ClusterRoleBindings().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}
	return crb.Items
}

func getServerVersion(clientset *kubernetes.Clientset) string {
	version, err := clientset.ServerVersion()
	if err != nil {
		log.Fatalln("failed to get services:", err)
	}
	return version.String()
}

// Exec

func execIntoPod(clientset *kubernetes.Clientset, pod *v1.Pod, command string, stdin io.Reader) (string, string, error) {
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := v1.AddToScheme(scheme); err != nil {
		return "", "", fmt.Errorf("error adding to scheme: %v", err)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&v1.PodExecOptions{
		Command: strings.Fields(command),
		Stdin:   stdin != nil,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(getConfig(), "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("error while creating Executor: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return "", "", fmt.Errorf("error in Stream: %v", err)
	}

	return stdout.String(), stderr.String(), nil

}

func testConnectionPodToService(pod v1.Pod, svc v1.Service) bool {

	clientset := getClientSet()

	// execute command
	var result bool
	command := "wget --spider --timeout=1 " + svc.Name + "." + svc.Namespace
	output, stderr, err := execIntoPod(clientset, &pod, command, nil)
	if len(stderr) != 0 {
		fmt.Println("STDERR:", stderr)
		result = false
	}
	if err != nil {
		fmt.Printf("Error occured while `exec`ing to the Pod %s, namespace %s, command %s\n", pod.Name, pod.Namespace, command)
		fmt.Println(err)
		result = false
	} else {
		fmt.Println("Output:")
		fmt.Println(output)
		result = true
	}

	return result
}

// func main() {
// 	var cs = getClientSet()
// 	var p = getPods(cs)
// 	var d = getDeployments(cs)
// 	var rs = getReplicaSets(cs)
// 	var ds = getDaemonSets(cs)
// 	var sts = getStatefulSets(cs)
// 	var s = getServices(cs)
// 	var scr = getSecrets(cs)
// 	var cm = getConfigMaps(cs)
// 	var sa = getServiceAccounts(cs)
// 	var ev = getEvents(cs)
// 	var ep = getEndpoints(cs)
// 	var pv = getPersistentVolumes(cs)
// 	var pvc = getPersistentVolumeClaims(cs)
// 	var r = getRoles(cs)
// 	var rb = getRoleBindings(cs)
// 	var cr = getClusterRoles(cs)
// 	var crb = getClusterRoleBindings(cs)
// 	var version = getServerVersion(cs)

// 	if len(p) > 0 {
// 		fmt.Println("pod: " + p[0].Name)
// 	}
// 	if len(d) > 0 {
// 		fmt.Println("deploy: " + d[0].Name)
// 	}
// 	if len(ds) > 0 {
// 		fmt.Println("ds: " + ds[0].Name)
// 	}
// 	if len(rs) > 0 {
// 		fmt.Println("rs: " + rs[0].Name)
// 	}
// 	if len(sts) > 0 {
// 		fmt.Println("sts: " + sts[0].Name)
// 	}
// 	if len(sts) > 0 {
// 		fmt.Println("sts: " + sts[0].Name)
// 	}
// 	if len(scr) > 0 {
// 		fmt.Println("scr: " + scr[0].Name)
// 	}
// 	if len(cm) > 0 {
// 		fmt.Println("cm: " + cm[0].Name)
// 	}
// 	if len(sa) > 0 {
// 		fmt.Println("sa: " + sa[0].Name)
// 	}
// 	if len(ev) > 0 {
// 		fmt.Println("ev: " + ev[0].Name)
// 	}
// 	if len(ep) > 0 {
// 		fmt.Println("ep: " + ep[0].Name)
// 	}
// 	if len(s) > 0 {
// 		fmt.Println("s: " + s[0].Name)
// 	}
// 	if len(pv) > 0 {
// 		fmt.Println("pv: " + pv[0].Name)
// 	}
// 	if len(pvc) > 0 {
// 		fmt.Println("pvc: " + pvc[0].Name)
// 	}
// 	if len(r) > 0 {
// 		fmt.Println("r: " + r[0].Name)
// 	}
// 	if len(rb) > 0 {
// 		fmt.Println("rb: " + rb[0].Name)
// 	}
// 	if len(cr) > 0 {
// 		fmt.Println("cr: " + cr[0].Name)
// 	}
// 	if len(crb) > 0 {
// 		fmt.Println("crb: " + crb[0].Name)
// 	}
// 	fmt.Println("version: " + version)

// 	command := "whoami"
// 	output, stderr, err := execIntoPod(cs, &p[0], command, nil)
// 	if len(stderr) != 0 {
// 		fmt.Println("STDERR:", stderr)
// 	}
// 	if err != nil {
// 		fmt.Printf("Error occured while `exec`ing to the Pod %q, namespace %q, command %q. Error: %+v\n", &p[0].Name, &p[0].Namespace, command, err)
// 	} else {
// 		fmt.Println("Output:")
// 		fmt.Println(output)
// 	}
// }
