# Minikube Setup
I hit a couple of issues in the setting up of Minikube, and needed to iterate through a couple to get something satisfactory.

### TL;DR: Minikube with VirtualBox
The TL;DR is that I used Minikube `v1.25.2` with _VirtualBox_ `v6.1.34 r150636 (Qt5.6.3)` as the driver.  
The kubernetes-cli versions are:
```shell
Client Version: v1.24.1
Kustomize Version: v4.5.4
Server Version: v1.23.3
```
I use `brew` to handle the install and update of all the packages except _VirtualBox_.
As I had an older setup of Minikube already present, I already had _VirtualBox_ installed, I just ran an internal update to grab the latest version. 

**Before starting Minikube**, run these commands to adjust how _VirtualBox_ handles dns resolution:
```shell
VBoxManage modifyvm "minikube" --natdnshostresolver1 off
VBoxManage modifyvm "minikube" --natdnsproxy1 on
```

Start minikube with:
```shell
 minikube start --driver=virtualbox
```

### Setup Debugging
If you are only interested in replicating the setup, then the information above should be sufficient and you do not need to read this subsection.
If you want to understand _why_ I ended up with this setup, I've included my debugging process and some thoughts along the way.

1. **Docker driver:**  
   Originally I opted to use _Docker_ as the driver, but hit issues trying to expose services use `NodePort`.
   I believe the issues are related to how Docker-Desktop-for-Mac handles networking with the host, and the [lack of the `docker0` bridge](https://docs.docker.com/desktop/mac/networking/#known-limitations-use-cases-and-workarounds).
   To debug this I created a deployment and service for a simple _nginx_ app, included in the submission under the `nginx_debugging` directory.  
   After deploying, I expect to be able to access the service at a path determined by Minikube using `minikube service nginx-svc --url`, however this never becomes available when using _Docker_ as the vm-driver.
   I didn't invest more time to find a workaround, partly because using other vm-drivers avoid the issue and partly because I wouldn't expect any workarounds to be portable outside of macOS.

1. **DNS Resolution with VirtualBox:**  
   I hit an issue when building the docker images when using Minikube with _VirtualBox_: the build fails with a DNS resolution error.
   Hunting this down, it seems to stem from a change in behaviour in how _busybox_ performs nslookups.
   This was introduced in `alpine3.13` (from which your images are built), but only for very specific circumstances... and it seems how I am using _VirtualBox_ is one of them.  
   Reading up the problem, I decided on [this workaround](https://github.com/alpinelinux/docker-alpine/issues/149#issuecomment-1110979984) which involces running the two `VBoxManage` commands before booting up the Minikube cluster.
   I believe this is changing how _VirtualBox_ decides to handle the DNS resolution, offloading the entire resolution to the host instead of just using the host system APIs.