#! /bin/sh

set -e

TAG=$1
DOCKERACC=$2

if [ -z "${TAG}" ]; then
	TAG="$(git rev-parse --short HEAD)"
	if [ -z "${TAG}" ]; then
		echo no TAG
		exit 2
	fi
	echo "No TAG specified. Using ${TAG}."
fi

if [ "${DOCKERACC}" = "" ]; then
	DOCKERACC=midonet
fi

for VARIANT in amd64-linux arm64v8-linux; do
	docker build -f Dockerfile-${VARIANT} -t ${DOCKERACC}/midonet-kube-controllers-${VARIANT}:${TAG} .
	docker build -f Dockerfile-node-${VARIANT} -t ${DOCKERACC}/midonet-kube-node-${VARIANT}:${TAG} .
done

echo "Now you can push images with the following commands:"
for VARIANT in amd64-linux arm64v8-linux; do
	echo "  docker push ${DOCKERACC}/midonet-kube-controllers-${VARIANT}:${TAG}"
	echo "  docker push ${DOCKERACC}/midonet-kube-node-${VARIANT}:${TAG}"
done
