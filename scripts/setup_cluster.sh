#!/usr/bin/env bash
DEV_USER=fusordev

oc cluster up

# Create and configure user
oc login -u system:admin
oc create user $DEV_USER
oadm policy add-cluster-role-to-user cluster-admin fusordev
oadm policy add-scc-to-user privileged fusordev
oadm policy add-scc-to-group anyuid system:authenticated
