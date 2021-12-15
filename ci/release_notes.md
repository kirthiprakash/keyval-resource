### Improvements

- Build the resource image based on `alpine:latest` instead of the possibly unstable `alpine:edge`.
- Improved the Concourse pipeline to re-build the resource image whenever either the Golang or Alpine base images are updated.
- Trigger Pull Request testing immediately using a GitHub Action whenever some PR is created or updated.
