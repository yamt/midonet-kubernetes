## How to build

### Prequisite

- docker
- [dep][dep]

[dep]: https://github.com/golang/dep

### Build docker images

<pre>
	% dep ensure
	% TAG=1.1
	% ./build.sh ${TAG}
</pre>

### Upload docker images

See the output of the above build.sh script.

### Create docker manifest lists

<pre>
	% TAG=1.1
	% ./createmanifest.sh ${TAG}
</pre>

Note: This uses ["docker manifest"][docker-manifest] command,
which is experimental.
You might need to enable the experimental feature on your docker
environment. And/or you might need to upgrade your docker.

[docker-manifest]: https://docs.docker.com/edge/engine/reference/commandline/manifest/

### Upload docker manifest lists

See the output of the above createmanifest.sh script.

