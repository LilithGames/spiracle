.PHONY: pre
pre:
	@go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.5.0

.PHONY: gen-crd
gen-crd: SHELL:=C:\Windows\System32\bash.exe
gen-crd:
	@bash `go list -f '{{ .Dir }}' -m k8s.io/code-generator@v0.21.1`/generate-groups.sh "deepcopy,client,informer,lister" github.com/LilithGames/spiracle/pkg/generated github.com/LilithGames/spiracle/pkg/apis "samplecontroller:v1alpha1" --output-base github.com/LilithGames/spiracle --go-header-file doc/boilerplate.go.txt

.PHONY: crd
crd: crd-object crd-manifests

.PHONY: crd-manifests
crd-manifests:
	@controller-gen crd:trivialVersions=true,preserveUnknownFields=false rbac:roleName=spiracle-role webhook paths=./... output:crd:artifacts:config=deploy/crd output:rbac:artifacts:config=deploy/rbac output:webhook:artifacts:config=deploy/webhook

.PHONY: crd-object
crd-object:
	@@controller-gen object:headerFile="doc/boilerplate.go.txt" paths="./..."

.PHONY: build
build:
	@GOOS=linux go build -o bin/ github.com/LilithGames/spiracle/...

.PHONY: run
run: build
	@wsl -e bin/spiracle

.PHONY: image image-ci
image: build
	@DOCKER_BUILDKIT=1 docker-compose -f deploy/build/docker-compose.yaml build spiracle

.PHONY: image-local
image-local:
	@docker-compose -f deploy/build/docker-compose.yaml build spiracle-local

.PHONY: push
push: crd build
	@REVISION=$$(git describe --tags) docker-compose -f deploy/build/docker-compose.yaml build spiracle-release
	@REVISION=$$(git describe --tags) docker-compose -f deploy/build/docker-compose.yaml push spiracle-release

.PHONY: install
install:
	@kubectl apply -k deploy
	@kubectl rollout restart statefulset.apps/spiracle -n default
	@kubectl rollout status statefulset.apps/spiracle -n default

.PHONY: deploy
deploy: image install

.PHONY: clean
clean:
	@kubectl delete -k deploy

.PHONY: install-sample
install-sample:
	@kubectl apply -f deploy/samples/roomingress.yaml

.PHONY: clean-sample
clean-sample:
	@kubectl delete -f deploy/samples/roomingress.yaml

.PHONY: tag
tag:
	@git tag $$(svu next)

.PHONY: release
release: tag push
	@git push --tags
