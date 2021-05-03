#!/bin/bash
. VERSION
git tag -d $VERSION
git push --delete origin $VERSION
git tag $VERSION
git push --tags
