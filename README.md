# github-to-terraform

## Description

`github-to-terraform` is used to retrieve your organization's GitHub
configuration and convert it into Terraform code. Currently supported functions:

- collaborators

## Installation

For OSX Homebrew:

```sh
brew tap shmileee/homebrew-tap
brew install github-to-terraform
```

## Prerequisites

Set `GITHUB_TOKEN` environmental variable:

```sh
export GITHUB_TOKEN=<gh pat>
```

## Usage

```sh
Get current GitHub configuration and create Terraform code for it

Usage:
  github-to-terraform [command]

Available Commands:
  collaborators Retrieve repository collaborators from GitHub and save them as
  Terraform resources
  completion    Generates bash completion scripts
  help          Help about any command
  version       Print the version

Flags:
  -h, --help   help for github-to-terraform

Use "github-to-terraform [command] --help" for more information about a command.
```

## Examples

Extract outside collaborators for all private repositories in organization:

```sh
github-to-terraform collaborators --org Appsilon
```

Extract collaborators from single public repository like this:

```sh
github-to-terraform collaborators --repo-type public --repo-name <repository>
```
