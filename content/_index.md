---
title: wp-git-clone
---

[![Build Status](https://ci.thegeeklab.de/api/badges/thegeeklab/wp-git-clone/status.svg)](https://ci.thegeeklab.de/repos/thegeeklab/wp-git-clone)
[![Docker Hub](https://img.shields.io/badge/dockerhub-latest-blue.svg?logo=docker&logoColor=white)](https://hub.docker.com/r/thegeeklab/wp-git-clone)
[![Quay.io](https://img.shields.io/badge/quay-latest-blue.svg?logo=docker&logoColor=white)](https://quay.io/repository/thegeeklab/wp-git-clone)
[![Go Report Card](https://goreportcard.com/badge/github.com/thegeeklab/wp-git-clone)](https://goreportcard.com/report/github.com/thegeeklab/wp-git-clone)
[![GitHub contributors](https://img.shields.io/github/contributors/thegeeklab/wp-git-clone)](https://github.com/thegeeklab/wp-git-clone/graphs/contributors)
[![Source: GitHub](https://img.shields.io/badge/source-github-blue.svg?logo=github&logoColor=white)](https://github.com/thegeeklab/wp-git-clone)
[![License: Apache-2.0](https://img.shields.io/github/license/thegeeklab/wp-git-clone)](https://github.com/thegeeklab/wp-git-clone/blob/main/LICENSE)

Woodpecker CI plugin to clone git repositories.

<!-- prettier-ignore-start -->
<!-- spellchecker-disable -->
{{< toc >}}
<!-- spellchecker-enable -->
<!-- prettier-ignore-end -->

## Usage

```YAML
clone:
  git:
    image: quay.io/thegeeklab/wp-git-clone
    settings:
      depth: 50
      lfs: false
```

### Parameters

<!-- prettier-ignore-start -->
<!-- spellchecker-disable -->
{{< propertylist name=wp-git-clone.data sort=name >}}
<!-- spellchecker-enable -->
<!-- prettier-ignore-end -->

## Build

Build the binary with the following command:

```Shell
make build
```

Build the container image with the following command:

```Shell
docker build --file Containerfile.multiarch --tag thegeeklab/wp-git-clone .
```

## Test

```Shell
docker run --rm \
  -e CI_REPO_CLONE_URL=https://github.com/octocat/Hello-World.git \
  -e CI_PIPELINE_EVENT=push \
  -e CI_COMMIT_SHA=553c2077f0edc3d5dc5d17262f6aa498e69d6f8e \
  -e CI_COMMIT_REF=refs/heads/master \
  -e CI_WORKSPACE=/tmp/wp_git_testrepo \
  -v $(pwd):/build:z \
  -w /build \
  quay.io/thegeeklab/wp-git-clone
```
