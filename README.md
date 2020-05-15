# boilr-runner

This is a [boilr template](https://github.com/tmrts/boilr) for creating a custom pipeline runners. Get started by cloning the project and installing the template:

```console
$ boilr template download drone/boilr-runner runner
```

create a project in directory my-runner:

```console
$ boilr template use runner my-runner
```

enter the runner kind (e.g. "pipeline")

```text
[?] Please choose a value for "Type" [default: "pipeline"]:
```

enter the runner type (i.e. "docker", "ssh", etc)

```text
[?] Please choose a value for "Type" [default: ""]:
```

enter the Docker image name:

```text
[?] Please choose a value for "Image" [default: "docker.io/octocat/hello-world"]:
```

enter the Go module name:

```text
[?] Please choose a value for "Module" [default: "github.com/octocat/hello-world"]:
```
