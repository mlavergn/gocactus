[![Build Status](https://github.com/mlavergn/gocactus/workflows/CI/badge.svg?branch=master)](https://github.com/mlavergn/gocactus/actions)
[![Go Report](https://goreportcard.com/badge/github.com/mlavergn/gocactus)](https://goreportcard.com/report/github.com/mlavergn/gocactus)
[![GoDoc](https://godoc.org/github.com/mlavergn/gocactus/src/rx?status.svg)](https://godoc.org/github.com/mlavergn/gocactus/src/rx)

# Go Cactus

Go Cactus is a module for OpenWRT devices to quickly toggle Cactus DNS servers

## Background

OpenWRT devices are excellent gateways for home internet service. When combined with private DNS services, they offer opportunities to prevent DNS based geofencing. However, for non-technical users, modifiying DNS settings is confusing and could result in a misconfigured router.

This module seeks to simplify the process by exposing a simple REST API that can be used to toggle private DNS via a simple web app.
