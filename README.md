## Welcome

We're really happy that you're considering joining us!
This challenge will help us understand your skills and will also be a starting point for the next interview.
We're not expecting everything to be done perfectly as we value your time but the more you share with us, the more we get to know about you!

This challenge is split into 3 parts:

1. Debugging
2. Implementation
3. Questions

If you find possible improvements to be done to this challenge please let us know in this readme and/or during the interview.

## The challenge

Pleo runs most of its infrastructure in Kubernetes.
It's a bunch of microservices talking to each other and performing various tasks like verifying card transactions, moving money around, paying invoices, ...
This challenge is similar but (a lot) smaller :D

In this repo, we provide you with:

- `invoice-app/`: An application that gets invoices from a DB, along with its minimal `deployment.yaml`
- `payment-provider/`: An application that pays invoices, along with its minimal `deployment.yaml`
- `Makefile`: A file to organize commands.
- `deploy.sh`: A file to script your solution
- `test.sh`: A file to perform tests against your solution.

### Set up the challenge env

1. Fork this repository
2. Create a new branch for you to work with.
3. Install any local K8s cluster (ex: Minikube) on your machine and document your setup, so we can run your solution.
> Documented setup of Minikube, plus some comments on decision-making for this setup, in the included `minikube_setup.md` file

### Part 1 - Fix the issue

The setup we provide has a :bug:. Find it and fix it! You'll know you have fixed it when the state of the pods in the namespace looks similar to this:

```
NAME                                READY   STATUS                       RESTARTS   AGE
invoice-app-jklmno6789-44cd1        1/1     Ready                        0          10m
invoice-app-jklmno6789-67cd5        1/1     Ready                        0          10m
invoice-app-jklmno6789-12cd3        1/1     Ready                        0          10m
payment-provider-abcdef1234-23b21   1/1     Ready                        0          10m
payment-provider-abcdef1234-11b28   1/1     Ready                        0          10m
payment-provider-abcdef1234-1ab25   1/1     Ready                        0          10m
```

#### Requirements

Write here about the :bug:, the fix, how you found it, and anything else you want to share about it.
> **The problem:**  
>The security context constraints require both apps to run as non-root, but the images are configured to run as root, and so Kubernetes does not launch the pods.  
> 
> **How I found it:**  
> Ostensibly the Dockerfile look fine, and builds from the outset.
> Checking the deployment YAML before launching, it also looks fine without any obvious errors (when considered alone).
> After launching, I probed the failing pods to reveal the error message `container has runAsNonRoot and image will run as root`.
> It is immediately obvious the problem is in the USER of the Dockerfile being root.  
> I've had projects building Helm charts for applications, and so spent time crafting suitable SCCs for the deployment and hitting similar problems before.
> 
> **The fix:**  
> There are two fixes: either remove the requirement for the pods to run as non-root, or adjust the images to not run as root.  
> I opted for the latter and adjusted the Dockerfiles, because generally it is better for containers not to be running as root.  


### Part 2 - Setup the apps

We would like these 2 apps, `invoice-app` and `payment-provider`, to run in a K8s cluster and this is where you come in!

#### Requirements

1. `invoice-app` must be reachable from outside the cluster.
2. `payment-provider` must be only reachable from inside the cluster.
3. Update existing `deployment.yaml` files to follow k8s best practices. Feel free to remove existing files, recreate them, and/or introduce different technologies. Follow best practices for any other resources you decide to create.
4. Provide a better way to pass the URL in `invoice-app/main.go` - it's hardcoded at the moment
5. Complete `deploy.sh` in order to automate all the steps needed to have both apps running in a K8s cluster.
6. Complete `test.sh` in order to perform tests against your solution and get successful results (all the invoices are paid) via `GET invoices`.


> I tried to make the set of commits logical in their progression, and the final version fulfils the requirements.
> I've added some comments on a couple of the points to try and express my thought-process a little more.
> 
> ### Comments on Q1
> I opted to use `NodePort`, as it was sufficient for the job in this toy setup.
> If the system was larger (ie: spanning multiple nodes) this wouldn't be sufficient and would need shifting to `LoadBalancer` at least.
> 
> ### Comments on Q4
> 
> For Q4, I opted to set the invoice-app to read the `HOST:PORT` from environment variables, which means this can be set at deployment time without needing to rebuild the image.
> I added logging of the full target URL to aid with debugging, but (depending on requirements) it might be worth catching situations where the URL is incorrectly set (or not set at all) and raising somewhere more obvious.
> I didn't add this in as it felt beyond the scope of the exercise.
> 
> ### Comments on Q6  
> Either by accident or design, the `test.sh` file was not included in the challenge repository.
> It wasn't a problem as I could back-engineer the simple interactions with the apps, however I don't know if there were specifics on how `test.sh` should look.
> For example, if the expectation was for `test.sh` to return a `pass|fail`.  
> My testing consisted of running `make pay` & `make invoice` and physically checking that no invoices remain.
> From the `README.md` description, I would guess there is a semi-complete script that calls the `/invoices/pay` API, and then checks that the returned dictionary is blank.
> Scripting this together is not a problem for me, just let me know if this is needed.
> Given free reign, I would do this using python (my preferred language atm).

### Part 3 - Questions

Feel free to express your thoughts and share your experiences with real-world examples you worked with in the past.

#### Requirements

1. What would you do to improve this setup and make it "production ready"?
2. There are 2 microservices that are maintained by 2 different teams. Each team should have access only to their service inside the cluster. How would you approach this?
3. How would you prevent other services running in the cluster to communicate to `payment-provider`?

> ## Part 3: Answers
> ### Q1:
> In no particular order, some improvements I would consider before moving to production...
> 
> - **improved testing:**  
>   My testing method is not automated, but even if it were, I am only testing the top-level functionality: I'm checking that everything works together to successfully pay all the invoices.
>   Following through the code, this invokes lots of individual pieces. For example:
>     - _invoice-app_ needs to successfully spin up the DB,
>     - the `getInvoices` function needs to successfully call the `payment-provider::/payments/pay` API
>     - the DB needs to correctly update the invoice status  
> 
>   A better more thorough approach would be to test these fundamental pieces in some isolated way, and building up to the top-level checks.
>   For simple applications, it doesn't seem so important but for large, complex stacks this can make debugging much easier and help to identify areas where testing is lacking.
>
> - **Helm Charts:**  
>   I would consider if a helm-chart can help simplify the deployment process here.
>   I've built Charts for deployment of some of our applications and find them much neater/efficient than deployment scripts.
>   However, I think the real advantage here is in linking configurations across different services.
>   For example, for _payment-provider_ we build a `service` and specify the `port` and `targetPort`.
>   The _invoice-app_ needs to know these to build the appropriate URL, if the (hard-coded) values across the two `deployment.yaml` files don't match, the whole process won't work, but not in an obvious way.
>   With Helm, these can be higher-level variables that are then populated across the YAML objects as needed.
>   It helps make a consistent environment, and it solidifies (in code) where the variables are referenced and how k8 objects interact with each other.  
>   On the other hand, these advantages are applicable when both services are contained in a single Helm Chart.
>   From the simplified setup this seems suitable. 
>   However if the two services were more complex and, for some reason, were better suited to be separate Helm Charts, it might not be so easy to link the configurations between the two.
>   
> - **Solidifying external access:**  
>   For the challenge purposes, I'm using Minikube (with only 1 node), and `NodePorts` to externalise access.
>   In production this wouldn't be suitable because it doesn't load-balance over more nodes, and (crucially) the access point to the service is ephemeral as it is based on the `nodeIP`.
>   Instead, an ingress + controller could be configured to provide a consistent access route to the service, regardless of where it is being hosted at that time.
> 
> - **pod-local database**  
>   I'm not sure if this is by design or by oversight, but the database holding the invoices lives inside the _invoice-app_, of which there are three replicas.
>   This means there are three independent sets of identical invoices, and calling `invoice-app::/invoices/pay` updates only one of these (chosen by the rules of selection that the service follows).  
>   This is almost certainly going to be a bad idea, because it's difficult (if not impossible) to specify which of the databases to interact with, meaning trying to track the correct status of invoices becomes very difficult.
>   Unless there is a good reason for otherwise, it is better to have the database exist as a standalone instance which any replica of _invoice-app_ can access and safely update information in.
> 
> - **logging:**  
>   I would consider adding more logging, particularly around failure.
>   Right now if the `pay` function within invoice-app `/invoices/pay` API call fails, there is no obvious notification to the user.
>   For example, if the `HOST:PORT` doesn't exist, the API call seems to fail silently (nothing in the invoice-app logs).
>   A notification in the logs of a failed call, or some handling of an unknown `HOST:PORT` would help with use and debugging.
> 
> ### Q2
> Depending on the extent of isolation, `namespaces` could provide a solution.
> If each microservice is assigned it's own namespace, members of the asssociated team can be granted permission to that namespace using RBAC rules.
> This should provide the isolation needed, but if these microservices need to interact with each other it might cause complications.  
> I've had services communicate across namespaces before, but I've not had to enforce isolation like the question describes.
> I'm not sure if maintaining isolation, but allowing particular services to communicate across namespaces, is possible with `namespaces` alone.
> It might require some additional setup (such as `NetworkPolicies`).
> 
> ### Q3
> In truth, I haven't had much experience with securing individual `Services`, because I've never had to do it.
> However, I believe `NetworkPolicies` can be used to achieve what the question asks.
> I don't know exactly how to implement the policies in a way to achieve the result, but I believe `NetworkPolicies` can allow the _payment-provider_ service to reject all traffic except that originating from the _invocie-app_ pods.
> 

---

## What matters to us?

Of course, we expect the solution to run, but we also want to know how you work and what matters to you as an engineer.
So, feel free to use any technology you want!
You can create new files, refactor, rename, ...

Ideally, we'd like to see your progression through commits, verbosity in your answers and all requirements met.
Don't forget to update the README.md to explain your thought process.
