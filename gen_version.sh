#!/bin/bash

ROOT="$(pwd)"

if BUILD_GIT_REVISION=$(git rev-parse HEAD 2> /dev/null); then if ! git diff-index --quiet HEAD; then BUILD_GIT_REVISION=${BUILD_GIT_REVISION}"-dirty"
    fi
else
    BUILD_GIT_REVISION=unknown
fi

# Check for local changes
if git diff-index --quiet HEAD --; then
  tree_status="Clean"
else
  tree_status="Modified"
fi

# Check for git branch and git dirty
BRANCH=$(git rev-parse --abbrev-ref HEAD)
GIT_DIRTY=$(git diff --no-ext-diff 2> /dev/null | wc -l)

# XXX This needs to be updated to accommodate tags added after building, rather than prior to builds
RELEASE_TAG=$(git describe)

# security wanted VERSION='unknown'
VERSION="${BUILD_GIT_REVISION}"
if [[ -n "${RELEASE_TAG}" ]]; then
  VERSION="${RELEASE_TAG}"
fi

# used by core/version
echo buildVersion       "${VERSION}"
echo buildGitRevision   "${BUILD_GIT_REVISION}"
echo buildUser          "$(whoami)"
echo buildHost          "$(hostname -f)"
echo buildStatus        "${tree_status}"
echo buildTime          "$(date +%Y-%m-%d--%T)"
echo buildBranch        "${BRANCH}"
echo buildGitDirty      "${GIT_DIRTY}"
