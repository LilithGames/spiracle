def expand(paths):
	return  [file for path in paths for file in listdir(path, recursive=True)]

registry = 'lilith-registry.cn-shanghai.cr.aliyuncs.com/avatar/hulucc/'
build_deps = ['api', 'client', 'cmd', 'config', 'controllers', 'infra', 'protocol', 'proxy', 'repos', 'services', 'utils']
build_ignores = ['**/*_test.go']
image_deps = ['./bin/spiracle']
deploy_deps = ['./deploy/']
deploy_image_deps = [registry+'spiracle:latest']
dockerfile = 'deploy/build/Dockerfile'
live_update = [sync('./bin/spiracle', '/spiracle/spiracle')]
entrypoint = ['/spiracle/spiracle', '-config', '/etc/spiracle/config.yaml']

allow_k8s_contexts('local')
load('ext://restart_process', 'docker_build_with_restart')

local_resource('build-linux', 'make build', deps=expand(build_deps), ignore=build_ignores)
docker_build_with_restart(registry+'spiracle', '.', entrypoint, dockerfile=dockerfile, live_update=live_update, only=image_deps)
# custom_build_with_restart(registry+'spiracle', 'make image-local', image_deps, entrypoint, live_update=live_update, disable_push=True)
k8s_yaml(kustomize('deploy'))
k8s_custom_deploy('deploy', 'make install', 'make clean', deploy_deps, image_deps=deploy_image_deps)
