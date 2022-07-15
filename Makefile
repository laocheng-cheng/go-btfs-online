# service runtime arguments
SHELL := /bin/bash
stage = $(shell tron-get-stage-from-branch --application-name btfs)
region = $(shell aws-eks-get-pod-region)
cluster = $(shell tron-get-cluster-name)
include stages/$(stage)
version = $(shell tron-get-tag-for-stage --stage-file stages/$(stage))

# do not change the variables below
service = status
repository := bt/$(service)
image := $(repository):$(version)
kubectl := kubectl -n $(namespace)
aws := aws --region $(region)

# ECR registry parameters (for login, push targets)
# This uses the current aws config profile (or $AWS_PROFILE)
account := $(shell $(aws) sts get-caller-identity --output text --query 'Account')
registry := $(shell echo $(account).dkr.ecr.us-east-1.amazonaws.com)
login := $(shell $(aws) ecr get-login --registry-ids $(account) --no-include-email)

default: build

build: login
	docker build -t $(image) .

# Logs into the docker registry
login:
	@echo $(registry)
	@$(login)

# Creates a new docker image repository (first time only)
repo: login
	$(aws) ecr create-repository --repository-name $(repository) || true
	$(aws) ecr set-repository-policy --repository-name $(repository) \
		--policy-text file://aws/ecr-repository-access-policy.json

immutable:
	$(aws) ecr put-image-tag-mutability --repository-name $(repository) \
		--image-tag-mutability IMMUTABLE

# Pushes the container into the remote registry.
push: login immutable
	docker tag $(image) $(registry)/$(image)
	$(aws) ecr describe-images --repository-name $(repository) --image-ids imageTag=$(version) \
		|| docker push $(registry)/$(image)

# Configure and deploy the service
git_sha := $(shell git rev-parse HEAD)
replacements="\
s/SERVICE_NAME/$(service)/g;\
s/VERSION/$(version)/g;\
s/ACCOUNT/$(account)/g;\
s/NAMESPACE/$(namespace)/g;\
s/HOSTNAME/$(hostname)/g;\
s/ENVIRONMENT/$(environment)/g;\
s/GIT_SHA/$(git_sha)/g;\
s/MIN_PODS/$(min_pods)/g;\
s/MAX_PODS/$(max_pods)/g;\
s/CPU_LIMITS/$(cpu_limits)/g;\
s/CPU_REQUESTS/$(cpu_requests)/g;\
s/MEMORY_LIMITS/$(memory_limits)/g;\
s/MEMORY_REQUESTS/$(memory_requests)/g;\
"

namespace:
	@cat kubernetes/namespace.yml | sed $(replacements) | $(kubectl) apply -f -

run: namespace
	@cat kubernetes/service-account.yml | sed $(replacements) | $(kubectl) apply -f -
	@cat kubernetes/ingress-nginx.yml | sed $(replacements) | $(kubectl) apply -f -
	@cat kubernetes/service.yml | sed $(replacements) | $(kubectl) apply -f -
	@cat kubernetes/hpa.yml | sed $(replacements) | $(kubectl) apply -f -
	@$(kubectl) rollout status deployments/$(service)-app

inject-db-secrets: namespace
	@tron-alter-k8s-secrets-btfs --service $(service) --append $(namespace) --namespace $(namespace)

# Stops the service on the remote Kubernetes cluster
# but does not remove the service.
stop:
	@$(kubectl) delete deployment $(service)-app --ignore-not-found
	@$(kubectl) delete service $(service)-service --ignore-not-found
	@$(kubectl) delete ingress $(service)-ingress --ignore-not-found

clean: stop
	@$(kubectl) delete service $(service)-service --ignore-not-found

purge: clean

watch:
	@$(kubectl) get pods -w

test: validate-cert
	grpcurl $(hostname):443 grpc.health.v1.Health/Check

validate-cert:
	@tron-get-cert-details-from-k8s-secret --secretname $(service)-tls --namespace $(namespace) --hostname $(hostname)

# Inject secrets from AWS to kube
inject-secrets: namespace
	@k8s-secrets-inject-secrets $(service) $(namespace) ./aws/secrets

# Inject secrets from file to aws
create-secrets:
	@aws-secretsmanager-create-secrets $(service) ./aws/secrets

# Inject secrets from file to aws
update-secrets:
	@aws-secretsmanager-update-secrets $(service) ./aws/secrets
