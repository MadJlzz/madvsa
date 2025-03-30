package main

import (
	"fmt"
	"strings"
)

func main() {
	destination := "file:///tmp/scans"
	d, ok := strings.CutPrefix(destination, "file://")
	fmt.Println(d, ok)
}

//func main() {
//	//cli, err := client.NewClientWithOpts(client.WithHost("unix:///run/user/1000/podman/podman.sock"), client.WithAPIVersionNegotiation())
//	cli, err := client.NewClientWithOpts(client.WithHost("unix:///run/docker.sock"), client.WithAPIVersionNegotiation())
//	if err != nil {
//		panic(err)
//	}
//
//	// List running containers
//	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
//	if err != nil {
//		panic(err)
//	}
//
//	for _, c := range containers {
//		fmt.Println("Container ID:", c.ID, "Image:", c.Image)
//	}
//}

// KO
//func main() {
//	fmt.Println("Welcome to the Podman Go bindings tutorial")
//
//	// Get Podman socket location
//	//sock_dir := os.Getenv("XDG_RUNTIME_DIR")
//	//socket := "unix:" + sock_dir + "/podman/podman.sock"
//
//	// Connect to Podman socket
//	connText, err := bindings.NewConnection(context.Background(), "unix://run/docker.sock")
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//
//	containerLatestList, err := containers.List(connText, &containers.ListOptions{})
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//
//	fmt.Println(containerLatestList)
//}

//func main() {
//
//	conn, err := net.Dial("unix", "/home/jklaer/workspace/repos/madvsa/trivy/aSocket.sock")
//	if err != nil {
//		panic(err)
//	}
//	defer conn.Close()
//
//	_, err = conn.Write([]byte("hello there"))
//	if err != nil {
//		panic(err)
//	}
//
//}

//func main() {
//
//	//conn, err := net.Dial("unix", "/run/docker.sock")
//	//if err != nil {
//	//	panic(err)
//	//}
//	//defer conn.Close()
//
//	cli := http.Client{
//		Transport: &http.Transport{
//			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
//				return net.Dial("unix", "/run/docker.sock")
//			},
//		}}
//
//	resp, err := cli.Get("http://localhost/containers/json")
//	if err != nil {
//		panic(err)
//	}
//	defer resp.Body.Close()
//
//	b, err := io.ReadAll(resp.Body)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(string(b))
//
//	//_, err = conn.Write([]byte("/v1.48/containers/json"))
//	//if err != nil {
//	//	panic(err)
//	//}
//
//}
