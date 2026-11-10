# Before you open a pull-request

## Linting
Every pull-request runs a linting and formatting check with ``golangci-lint v2``, (the exact version can be seen in the [prek config file](prek.toml)). Please ensure that your branch passes this check *before*
submitting a pull-request.

### Linting with prek (reccomended)
This project uses [prek](https://github.com/j178/prek) to setup and manage commit-hooks and linting, and it is highly recommended that you install it before working with this repository.
After installation, clone the project and run:
```bash
prek install
```

Now a ``golangci-lint`` check that aligns with the CI-workflow will run as a pre-commit hook.
You can also manually run the analysis on staged files using:

```bash
prek run
```
