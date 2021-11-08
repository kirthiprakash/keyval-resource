
[![Docker Stars](https://img.shields.io/docker/stars/gstack/keyval-resource.svg?style=plastic)](https://registry.hub.docker.com/v2/repositories/gstack/keyval-resource/stars/count/)
[![Docker pulls](https://img.shields.io/docker/pulls/gstack/keyval-resource.svg?style=plastic)](https://registry.hub.docker.com/v2/repositories/gstack/keyval-resource)
<!--
[![Concourse Build](https://ci.gstack.io/api/v1/teams/gk-plat-devs/pipelines/keyval-resource/jobs/build/badge)](https://ci.gstack.io/teams/gk-plat-devs/pipelines/keyval-resource)
-->
[![dockeri.co](https://dockeri.co/image/gstack/keyval-resource)](https://hub.docker.com/r/gstack/keyval-resource/)

# Concourse CI Key Value Resource

Implements a resource that passes key values between jobs without using any
external resource such as Git/S3 etc.

The key/value pairs are serialized in the `version` JSON objects, stored in
the Concourse SQL database. As such, they are designed to hold small textual,
non-secret configuration data.

In terms of pipeline design, secrets are supposed to be stored in a vault like
CredHub instead, and binaries or large text files are supposed to be stored
on more relevant persistent storage like Git or S3.

## Thanks

This resource is a fork of the [keyval resource][moredhel_gh] by @moredhel.

Compared to the [original `keyval` resource][swce_gh] from @SWCE, writing
key/value pairs as plain files in some resource folder is more consistent
with usual conventions in Concourse, when it comes to storing anything in
step artifacts.

Writing/reading files is always easier in Bash scripts than parsing some Java
Properties file. Much less boilerplate code is required.

[moredhel_gh]: https://github.com/moredhel/keyval-resource
[swce_gh]: https://github.com/SWCE/keyval-resource

## Source Configuration

``` YAML
resource_types:
  - name: key-value
    type: registry-image
    source:
      repository: gstack/keyval-resource
      
resources:
  - name: key-value
    type: key-value
```

#### Parameters

*None.*

## Behavior

### `check`: Report the latest stored key-value pairs

This is a version-less resource so `check` behavior is no-op.

It will detect the latest store key/value pairs, if any, and won't provide any
version history.

#### Parameters

*None.*

### `in`: Fetch the latest stored key-value pairs from the Concourse SQL database

Fetches the given key & values from the stored resource version JSON (in the
Concourse SQL database) and write them in their respective files where the
key is the file name and the value is the file contents.

```yaml
version:
    my_secret: secret_value
```

would result in:

```sh
$ cat resource/my_secret
secret_value
```

#### Parameters

*None.*

### `out`: Store new set of key-value pairs to the Concourse SQL database

Converts each file in the artifact directory designated by `directory` to a
set of key-value pairs, where file names are the keys and file contents are
the values. This set of key-value pairs is persisted in the `version` JSON
object, to be stored in the Concourse SQL database.

A value from a file in `directory` can be overridden by a matching key with
different value in the dictionary given as the `overrides` parameter. If you
need to store some Concourse `((vars))` value in a key-value resource, then
add it to the `overrides` parameter of some `put` step.

#### Parameters

- `directory`: *Required.* The artifact directory to be scanned for files, in
  order to generate key-value pairs

- `overrides`: *Optional.* A dictionary of key-value pairs that will override
  any matching pair with same key found in `directory`.


## Examples

```yaml
resource_types:
  - name: key-value
    type: registry-image
    source:
      repository: gstack/keyval-resource

resources:
  - name: build-info
    type: key-value

jobs:

  - name: build
    plan:
      - task: build
        file: tools/tasks/build/task.yml # <- must declare a 'build-info' output artifact
      - put: build-info
        params:
          directory: build-info

  - name: test-deploy
    plan:
      - in_parallel:
          - get: build-info
            passed: [ build ]
      - task: test-deploy
        file: tools/tasks/task.yml # <- must declare a 'build-info' input artifact
```

The `build` task writes all the key-value pairs it needs to pass along in
files inside the `build-info` output artifact directory.

The `test-deploy` job then reads the files from the `build-info` resource,
which produces a `build-info` artifact directory to be used by the
`test-deploy` task.

## Development

### Running the tests

Golang unit tests can be run from some shell command-line with Ginkgo, that
has [to be installed](https://github.com/onsi/ginkgo#getting-started) first.

```bash
make test
```

These unit test are embedded in the `Dockerfile`, ensuring they are
consistently run in a determined Docker image providing proper test
environment. Whenever the tests fail the Docker build will be stopped.

In order to build the image and run the unit tests, use `docker build` as
follows:

```bash
docker build -t keyval-resource .
```

### Contributing

Please make all pull requests to the `master` branch and ensure tests pass
locally.
