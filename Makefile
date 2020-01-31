# Get current directory
DIR := ${CURDIR}

.PHONY: setup
setup:
	@which ./bin/openapi-gen > /dev/null || go build -o ./bin/openapi-gen k8s.io/kube-openapi/cmd/openapi-gen

.PHONY: generate
generate: setup
	@operator-sdk generate k8s
	@operator-sdk generate crds
	@./bin/openapi-gen --logtostderr=true \
	    -o "" -i ./pkg/apis/credstash/v1alpha1 \
	    -O zz_generated.openapi \
	    -p ./pkg/apis/credstash/v1alpha1 \
	    -h ./hack/boilerplate.go.txt -r "-"

.PHONY: semantic-release
semantic-release:
	@npm ci
	@npx semantic-release

.PHONY: semantic-release-dry-run
semantic-release-dry-run:
	@npm ci
	@npx semantic-release -d

package-lock.json: package.json
	@npm install