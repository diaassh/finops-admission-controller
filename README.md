# FinOps Kubernetes Admission Controller

A Go-based Kubernetes admission webhook that estimates the monthly cost of application workloads and rejects Deployments that exceed their configured budgets.

## Problem Statement

Kubernetes workloads consume CPU and memory resources, which contribute to infrastructure costs. Without budget controls, applications can be deployed with resource requirements that exceed the intended spending limit.

This project helps prevent that by checking a workload's estimated monthly cost before Kubernetes accepts a Deployment.

## Features

- Calculates estimated monthly CPU and memory costs.
- Supports multiple application replicas.
- Validates resource requests and budget values.
- Accepts or rejects workloads based on their monthly budget.
- Integrates with Kubernetes using a Validating Admission Webhook.
- Uses HTTPS with TLS certificates.
- Includes automated Go tests.

## Tech Stack

- **Language:** Go
- **Containerization:** Docker
- **Orchestration:** Kubernetes
- **Local Kubernetes cluster:** kind
- **Kubernetes libraries:** `k8s.io/api` and `k8s.io/apimachinery`

## How It Works

1. A user submits a Deployment to Kubernetes.
2. Kubernetes sends an admission request to the FinOps webhook.
3. The webhook reads the CPU and memory requests, replica count, and monthly budget.
4. The cost calculator estimates the monthly workload cost.
5. If the estimated cost exceeds the budget, the webhook rejects the Deployment.
6. Otherwise, Kubernetes is allowed to continue creating the Deployment.

## Cost Calculation

The project uses illustrative pricing assumptions:

- CPU: ₹2 per core-hour
- Memory: ₹0.50 per GiB-hour
- Monthly hours: 730

**Formula:**

Monthly cost = (CPU cores × CPU price + Memory GiB × Memory price) × 730 × Replicas

For example, a workload requesting 2 CPU cores and 4 GiB of memory across 3 replicas has an estimated monthly cost of ₹13,140.

These prices are examples for demonstrating the calculation. They are not actual cloud-provider prices.

## Budget Enforcement Example

| Estimated monthly cost | Monthly budget | Decision |
|---:|---:|---|
| ₹13,140 | ₹10,000 | Rejected |
| ₹383.25 | ₹15,000 | Allowed |

The project was tested against a local Kubernetes cluster. Kubernetes rejected the over-budget test Deployment and accepted the within-budget test Deployment.

## Run the Tests

```powershell
go test .
```

## Build the Application

```powershell
go build .
```

## Build the Docker Image

```powershell
docker build -t finops-webhook:local .
```

## Deploy to a Local kind Cluster

The project uses Kubernetes manifests to create the TLS Secret, webhook Deployment, Service, and ValidatingWebhookConfiguration.

The Docker image must be loaded into the cluster before deploying:

```powershell
kind load docker-image finops-webhook:local --name finops-cluster
```

The TLS Secret and webhook configuration must also be created before testing admission requests.

## Test Budget Enforcement

Submit an over-budget workload:

```powershell
kubectl apply -f .\test-over-budget.yaml
```

The expected result is a rejection message explaining that the estimated cost exceeds the configured budget.

Submit a within-budget workload:

```powershell
kubectl apply -f .\test-within-budget.yaml
```

The expected result is successful Deployment creation, provided the cluster can schedule and run the workload.

## Future Improvements

- Use separate Certificate Authority and server certificates with automated certificate management.
- Improve validation for additional Kubernetes workload types and resource configurations.
- Support configurable cloud pricing models.
- Add structured logging and monitoring.
- Add configurable budgets for teams and namespaces.
- Add integration tests for more edge cases.

## Disclaimer

This is a learning and portfolio project. The cost model uses simplified example prices and is not intended to predict actual cloud bills. Production use would require stronger certificate management, security hardening, and more comprehensive resource validation.