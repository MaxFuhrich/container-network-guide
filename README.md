# How to make Docker containers communicate within a custom bridge network
<!-- TOC -->

- [How to make Docker containers communicate within a custom bridge network](#how-to-make-docker-containers-communicate-within-a-custom-bridge-network)
    - [Prerequisites](#prerequisites)
    - [Getting the images](#getting-the-images)
    - [Creating a custom bridge network](#creating-a-custom-bridge-network)
    - [Connect the containers to the network](#connect-the-containers-to-the-network)
    - [Wrap Up](#wrap-up)

<!-- /TOC -->
Docker is an open platform that allows an application to be run in an isolated environment. These isolated applications are called containers and they are playing an important role in modern cloud technologies. While Docker images (which are used to create containers) can be run inside platforms like Kubernetes and Cloud Foundry, it is also possible to use only Docker to create and run containers in their own environment and network. Docker uses containerd as container runtime which can also be used by Kubernetes up to version [1.20](https://kubernetes.io/blog/2020/12/02/dont-panic-kubernetes-and-docker/ "Kubernetes deprecating Docker as runtime").
Therefore it is easy for developers to create applications and images with Docker that can be used by different platforms.

If you are running your applications in Docker, situations can occur in which a container in an isolated network is needed (e.g. if the database must not be exposed to outer networks but still be accessible in a way). This can be achieved by creating a user-defined bridge network which is a link layer between different network segments and connecting containers to the network. A bridge network in Docker is a software bridge that allows containers that are connected to it to communicate with each other while being isolated from containers that aren't connected to this network. If you want to learn more about bridge networking with Docker you can have a look at their [docs](https://docs.docker.com/network/bridge/ "Docker bridge networks").

This tutorial describes step-by-step how you can set up a container (MongoDB) in a custom bridge network and communicate with it via another container (a simple HTTP API) which is in the same network and how you can connect to this container (and make some calls to the API) by exposing a port.
![alt text](container-tutorial.jpg "Overview")
*Overview of the custom bridge network and its containers*

## Prerequisites
* Docker
* Enough disk space, which is about 500 MB for the MongoDB image (the image size of MongoDB may change, though) and 20 MB for the Image made for this tutorial

You won't need anything else than Docker for this tutorial, since there is already an image of the application that is going to be containerized. If you have to learn the basics of Docker first or want to create your own image, you can find some great guides and tutorials on their [website](https://docs.docker.com/get-started/ "Docker get started").

## Getting the images
First, we get the official MongoDB image from the [Docker Hub](https://hub.docker.com/_/mongo "MongoDB Docker Hub"). This can be done by using the CLI command **docker pull** followed by the name of the image which is **mongo** in our case (this will pull the MongoDB image with the tag **latest**). The complete command looks like this:
```
docker pull mongo
```
Then, we get the image that was created for this tutorial:
```
docker pull maxfuhrich/container-network-example
```
If you prefer to build your own image instead of using the example image, feel free to use the code (and the Dockerfile) from [this github repository](https://github.com/MaxFuhrich/container-network-guide "Container Network Tutorial").
To see if the images have been pulled successfully we can list them with:
```
docker image ls
```
This will show all images that are available locally. If everything worked, we should find both images we just pulled. If you have pulled images before, there may be listed more than the two we need.
## Creating a custom bridge network
Next, we want to create a custom network. This can be done by using the CLI command **docker network create** followed by the name the network should have which is **hello-network** in our case. The complete command looks like this:
```
docker network create hello-network
```
This will create a bridge network with the name **hello-network**. A [bridge network](https://docs.docker.com/network/bridge/ "Use bridge networks") is a network that lets containers communicate with each other while being isolated from containers that are not connected to the network.

To check if the network has been successfully created we can use the command:
```
docker network ls
```
If our network is listed there, it has been created!
```
docker network ls
NETWORK ID     NAME            DRIVER    SCOPE
e99d15db1b9a   bridge          bridge    local
cc423edba456   hello-network   bridge    local
c5096c5db5a1   host            host      local
ddf6427be1bc   none            null      local
```
By default, the used driver for the network is the bridge driver. An overview of the different networking drivers and their use cases can be found [here](https://www.docker.com/blog/understanding-docker-networking-drivers-use-cases/ "Docker networking drivers use cases").

**Note**: If you don't need your network anymore after this tutorial, you can remove it with:
```
docker network rm hello-network
```
## Connect the containers to the network
Now, let's create and run the MongoDB container and connect it to the network **hello-network**. The container can be run with the command **docker run** followed by the name of the image, which is **mongo** in our case:
```
docker run -d --name mongodb --network hello-network mongo
```
**-d** runs the container detached so that we can use the terminal again.

**--name** sets the name of the container we create (**mongodb**). If **--name** is not defined, a random string will be assigned as container name.

**--network** sets the network, that this container will be connected to (**hello-network**).

**Important**: Don't change the name of the MongoDB container as the other container uses this name to resolve it into an IP-address. This is called *automatic service discovery*. If you don't want to use automatic service discovery, you can use the IPv4Address of the MongoDB container that you get if you inspect the network (this won't work with the code as is, because the HTTP container tries to connect to the host **mongodb**).

Now, the only thing that is left is to create and run the HTTP container, connect it to **hello-network** and expose a port so that the outside can access it:
```
docker run -d --name http-container --network hello-network -p 8080:8080 maxfuhrich/container-network-example
```
**-p** "host-port:container-port" publishes/maps a TCP port in the container (8080) on a port on the host (8080). If the port 8080 on the host is already in use, you can change the port to any other one that is free.

**Note**: You can omit the **-d** flag if you want to see the output of the HTTP application in your terminal.

To check if our containers are running in our custom bridge network we can inspect it with the command **docker inspect** followed by the name of the network, **hello-network**:
```
docker inspect hello-network
```
Which outputs:
```
[
    {
        "Name": "hello-network",
        "Id": "cc423edba4561250791ecc63970cccdb4341b644864b27443dcc1ff5fbdacc2a",
        "Created": "2021-01-12T14:50:53.9820995Z",
        "Scope": "local",
        "Driver": "bridge",
        "EnableIPv6": false,
        "IPAM": {
            "Driver": "default",
            "Options": {},
            "Config": [
                {
                    "Subnet": "172.18.0.0/16",
                    "Gateway": "172.18.0.1"
                }
            ]
        },
        "Internal": false,
        "Attachable": false,
        "Ingress": false,
        "ConfigFrom": {
            "Network": ""
        },
        "ConfigOnly": false,
        "Containers": {
            "48ca87014d0f6b656fbda0ffd8c6d1402cfab0dda0dd1032abcf45a1bb64b990": {
                "Name": "http-container",
                "EndpointID": "...",
                "MacAddress": "00:00:aa:00:00:00",
                "IPv4Address": "172.18.0.3/16",
                "IPv6Address": ""
            },
            "83cd03b01c5c37295d8861a42cffd7aa3fbba4856c9b158c3e10866008b6a902": {
                "Name": "mongodb",
                "EndpointID": "...",
                "MacAddress": "00:00:aa:00:00:01",
                "IPv4Address": "172.18.0.2/16",
                "IPv6Address": ""
            }
        },
        "Options": {},
        "Labels": {}
    }
]
```
If everything went right we should be able to call *localhost:8080/hello* (*localhost:8080* works too) in our browser which returns a string.

The endpoint */add* creates a new entry that only contains a string with the current time in our database (and returns this entry to the caller):

![alt text](tutorial-add.jpg "Endpoint /add")
**Note**: In practice this endpoint should use POST but for this guide GET is used so that it can be accessed via browser.

The endpoint */history* shows all elements of the database:

![alt text](tutorial-history.jpg "Endpoint /history")

Congratulations, we have successfully created a multi-container application in a custom bridge network!

Since this is the end of the tutorial we should do some cleanup. First, we have to stop the containers. This can be done by using the command **docker container stop** followed by the names of the containers we want to stop, **http-container** and **mongodb**. Next, we can remove the containers with the command **docker container rm** followed by the container names **http-container** and **mongodb**. Finally we remove our custom network with **docker network rm** and the name of the network **hello-network**.
In our CLI, the three commands will look like this:
```
docker container stop http-container mongodb
docker container rm http-container mongodb
docker network rm hello-network
```

## Wrap Up
In this guide you learned how to create a custom network, connect containers to it, let them communicate via automatic service discovery, and expose a port of one container to the host. This approach also works with containers of other applications but keep in mind which IP/Host the API Container tries to connect to.

Further reading:
* [Docker Docs: Networking with standalone containers](https://docs.docker.com/network/network-tutorial-standalone/ "Docker network tutorial")
