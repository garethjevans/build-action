# Build Action

[![build-and-test](https://github.com/garethjevans/build-action/actions/workflows/build-and-test.yaml/badge.svg)](https://github.com/garethjevans/build-action/actions/workflows/build-and-test.yaml)
[![golangci-lint](https://github.com/garethjevans/build-action/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/garethjevans/build-action/actions/workflows/golangci-lint.yml)
[![Build and Publish](https://github.com/garethjevans/build-action/actions/workflows/publish-image.yaml/badge.svg)](https://github.com/garethjevans/build-action/actions/workflows/publish-image.yaml)

TODO

## Usage

### Auth

  - `server`: Host of the API Server.
   
  - `ca-cert`: CA Certificate of the API Server.

  - `token`: Service Account token to access kubernetes.

  - `namespace`: _(required)_ The namespace to create the build resource in.

### Image Configuration

  - `destination`: _(required)_

  - `env`: 

  - `serviceAccountName`: Name of the service account in the namespace, defaults to `default`  

### Basic Configuration

```yaml
- name: Build Image
  id: build
  uses: garethjevans/build-action@main
  with:
    # auth
    server: ${{ secrets.SERVER }}
    token: ${{ secrets.TOKEN }}
    ca_cert: ${{ secrets.CA_CERT }}
    namespace: ${{ secrets.NAMESPACE }}
    # image config
    destination: gcr.io/project-id/name-for-image
    env: |
      BP_JAVA_VERSION=17
```

### Outputs

  - `name`: The full name, including sha of the built image.

### Example

```yaml
- name: Do something with image
  run:
    echo "${{ steps.build.outputs.name }}"
```

## License

TODO The scripts and documentation in this project are released under the [MIT License](LICENSE).

## Contributions

TODO Contributions are welcome! See [Contributor's Guide](docs/contributors.md)
