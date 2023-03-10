# mu-app

![Go Report Card](https://goreportcard.com/badge/github.com/neuralnorthwest/mu-app)
![GitHub](https://img.shields.io/github/license/neuralnorthwest/mu-app?style=plastic)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/neuralnorthwest/mu-app?style=plastic)
![GitHub Workflow Status (with branch)](https://img.shields.io/github/actions/workflow/status/neuralnorthwest/mu-app/cicd.yaml?branch=develop&style=plastic)
![GitHub search hit counter](https://img.shields.io/github/search/neuralnorthwest/mu-app/goto?style=plastic)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/neuralnorthwest/mu-app?style=plastic)
![Lines of code](https://img.shields.io/badge/lines%20of%20code-1k-blue?style=plastic)
![Status](https://img.shields.io/badge/status-in%20development-orange?style=plastic)

This is a complete sample of all Mu functionality. It is a good starting point
for a new project based on Mu.

## Creating a new app

To create a new app based on `mu-app`, you can use the following command:

```bash
go run github.com/neuralnorthwest/mu-app/cmd/new@latest my-app
```

This takes care of cloning the `mu-app` repository, removing the `.git` folder,
and renaming all references to `mu-app` to `my-app`.

### Enabling badges

By default, the badges at the top of this `README.md` will be removed. To
keep them, you need to specify the GitHub owner and repository name that you
will be pushing to. Do with with the `--target` flag:

```bash
go run github.com/neuralnorthwest/mu-app/cmd/new@latest my-app --target my-name/my-app
```

## Developer quick start

If you want work on `mu-app`, you can use the following commands to get started:

```bash
git clone https://github.com/neuralnorthwest/mu-app.git
cd mu-app
make setup-dev
```

### Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for more information.

## Code of Conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for more information.
