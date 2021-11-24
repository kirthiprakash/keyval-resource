### Features

- Properly implements the key-value resource contract, using file names as keys and file contents as values, thus conforming the `ConfigMap` format from Kubernetes, and conventions established by nearly all Concourse resources
- Fixed `put` steps
- Unit tests passing, and properly blocking the Docker build whenever failing
- Concourse build pipeline for faster publishing of new versions
- Requires the `directory` parameter for `put` steps
