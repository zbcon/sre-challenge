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

## What matters to us?

Of course, we expect the solution to run, but we also want to know how you work and what matters to you as an engineer.
So, feel free to use any technology you want!
You can create new files, refactor, rename, ...

Ideally, we'd like to see your progression through commits, verbosity in your answers and all requirements met.
Don't forget to update the README.md to explain your thought process.
