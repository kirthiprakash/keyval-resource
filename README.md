
[![Docker Stars](https://img.shields.io/docker/stars/gstack/keyval-resource.svg?style=plastic)](https://registry.hub.docker.com/v2/repositories/gstack/keyval-resource/stars/count/)
[![Docker pulls](https://img.shields.io/docker/pulls/gstack/keyval-resource.svg?style=plastic)](https://registry.hub.docker.com/v2/repositories/gstack/keyval-resource)
[![Docker build status](https://img.shields.io/docker/build/gstack/keyval-resource.svg)](https://github.com/gstackio/keyval-resource)
[![Docker Automated build](https://img.shields.io/docker/automated/gstack/keyval-resource.svg)](https://github.com/gstackio/keyval-resource)

[![dockeri.co](http://dockeri.co/image/gstack/keyval-resource)](https://hub.docker.com/r/gstack/keyval-resource/)

# Concourse CI Key Value Resource

Implements a resource that passes key values between jobs without using any
external resource such as Git/S3 etc.

The key/value pairs are serialized in the `version` JSON objects, stored in
the Concourse SQL database. As such, they are desinged to hold small textual,
non-secret configuration data.

In terms of pipeline design, secrets are supposed to be stored in a vault like
CredHub instead, and binaries or large text files are supposed to be sotred
on more relevant persistent storage like Git or S3.

## Thanks

This resource is a fork of the [keyval resource][moredhel_gh] by @moredhel.

Compared to the [original `keyval` resource][swce_gh] from @SWCE, writing
key/value pairs as plain files in some resource folder is more consistent
with usual conventions in Concourse, when it comes to storing anything in
step artifacts.

Writing/reading files is always easier in Bash scripts than parsing some Java
Properties file. much less biolerplate code is required.

[moredhel_gh]: https://github.com/moredhel/keyval-resource
[swce_gh]: https://github.com/SWCE/keyval-resource

## Source Configuration

``` YAML
resource_types:
  - name: keyval
    type: docker-image
    source:
      repository: gstack/keyval-resource
      
resources:
  - name: keyval
    type: keyval
```

#### Parameters

*None.*

## Behavior

### `check`: Produce a single dummy key

This is a version-less resource so `check` behavior is no-op.

It will detect the latest store key/value pairs, if any, and won't provide any
version history.

#### Parameters

*None.*

### `in`: Report the given time.

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

- `directory`: *Required.* The artifact directory to be scanned for files, in
  order to generate key-value pairs

- `overrides`: *Optional.* A dictionary of key-value pairs that will override
  any key-value pairs found in `directory`.


### `out`: Consumes the given properties file

Does the reverse of `in`, where each file in the resource directory is
converted to a key (the file name) and a value (the file contents) that is
added to the `version` JSON object, to be stored in the Concourse SQL
database.

#### Parameters

*None.*

## Examples

```yaml
resource_types:
  - name: keyval
    type: registry-image
    source:
      repository: gstack/keyval-resource

resources:
  - name: keyval
    type: keyval

jobs:

  - name: build
    plan:
      - task: build
        file: tools/tasks/build/task.yml
      - put: keyval
        params:
          directory: build-info

  - name: test-deploy
    plan:
      - in_parallel:
          - get: keyval
            passed: [ build ]
      - task: test-deploy
        file: tools/tasks/task.yml
```

The build job writes all the key values it needs to pass along in files inside
the `build-info` directory. The `test-deploy` job then reads the files in the
`keyval` directory and use them as necessary.

## Development

### Prerequisites

* golang is *required* - version 1.9.x is tested; earlier versions may also
  work.
* docker is *required* - version 17.06.x is tested; earlier versions may also
  work.
* godep is used for dependency management of the golang packages.

### Running the tests

**NOTE**: Tests have not yet been rewritten to reflect the updated configuration

The tests have been embedded with the `Dockerfile`; ensuring that the testing
environment is consistent across any `docker` enabled platform. When the docker
image builds, the test are run inside the docker container, on failure they
will stop the build.

Run the tests with the following command:

```sh
docker build -t keyval-resource .
```

### Contributing

Please make all pull requests to the `master` branch and ensure tests pass
locally.
