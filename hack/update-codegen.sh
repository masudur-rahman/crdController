#!/usr/bin/env bash

vendor/k8s.io/code-generator/generate-groups.sh all \
	github.com/masudur-rahman/crdController/pkg/client \
	github.com/masudur-rahman/crdController/pkg/apis \
	controller.crd:v1beta1

