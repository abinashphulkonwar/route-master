# Route-Master

Route-Master is a lightweight, open-source load balancer and reverse proxy tool built in Go. It's designed to distribute incoming network traffic across multiple servers to improve performance, reliability, and scalability.

## Features

- Load balancing: Distribute incoming network traffic across multiple servers to optimize resource usage and avoid overloading any single server.
- Reverse proxy: Serve as a gateway between clients and servers to protect the identity of the server and improve security.
- Dynamic configuration: Update the routing rules and target servers in real-time, without requiring a restart or downtime.
- HTTPS support: Route-Master supports HTTPS traffic, enabling secure communication between clients and servers.
- Logging: Track and monitor incoming requests, server responses, and errors with detailed logs.

# Installation

To install Route-Master, you need to have Go installed on your system. You can install it by following the instructions on the official website: https://golang.org/doc/install.

Once you have Go installed, you can download and install Route-Master using the following command:

```sh
$ go get github.com/abinashphulkonwar/route-master
```
## build binary for your operation system
```
$ go build github.com/abinashphulkonwar/route-master
```
or download from [![LinkedIn]](https://www.linkedin.com/in/abinash-phulkonwar-775b521a5/)
https://github.com/abinashphulkonwar/route-master/releases/download/v1.0.0

# Configuration

Route-Master uses a YAML configuration file to specify the routing rules and target servers. You can define multiple routes and their corresponding target servers in the configuration file. Here's an example configuration file:

```yaml
server:
  host: "localhost"
  port: 3002

node:
  - node1:
    scheme: "http"
    target:
      - "server1.8080"
      - "server2.8081"
    path: "/api"
  - node2:
    scheme: "http"
    target:
      - "server3:8443"
      - "server4:8080"
    path: "/auth"
```

This configuration file defines two routes: one for the /api path and another for the /auth path. The /api route targets server1 and server2, server1 listening on port 8080 and server2 listening on port 8081. The /auth route targets server3 and server4, server3 listening on port 8443 and server4 listening on port 8080

You can also specify additional settings such as load balancing method, health checks, and timeouts in the configuration file. For more information on configuring Route-Master, please refer to the documentation.

# Usage

To start Route-Master, run the following command:

```sh
$ route-master -f config.yaml
```

This will start Route-Master with the specified configuration file. You can then send requests to Route-Master using the defined routes and it will distribute the traffic across the target servers.

# License

Route-Master is released under the MIT License. See the LICENSE file for more information.

# Contributing

We welcome contributions to Route-Master! To contribute, please create a pull request with your changes. For major changes, please open an issue first to discuss what you would like to change.

# Contact

If you have any questions or issues, please feel free to contact us at abinashphulkonwar98@gmail.com

[![LinkedIn](https://img.shields.io/badge/-LinkedIn-blue?style=flat-square&logo=Linkedin&logoColor=white&link=https://www.linkedin.com/in/abinash-phulkonwar-775b521a5/)](https://www.linkedin.com/in/abinash-phulkonwar-775b521a5/)

Thank you for using Route-Master!
