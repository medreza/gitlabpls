# gitlabpls

GitLab, please! A CLI tool to help you run new GitLab CI/CD pipeline.
It generates new GitLab pipeline URL with configurable set of pipeline variables, so you don't have to type them again 
 and again.

## How to configure and install

1. Clone this repo

2. Head to `gitlabpls.config.toml` for configs. Place your pipeline variables under `vars.<your own vars key>` and 
the app will turn them into commands :sparkles:

    For example:
    ```
    [vars.myvariables]
    ENV = "dev"
    GIT_BRANCH = "${GIT_HEAD}"
    ```
    The variables defined within `vars.myvariables` can be used with `gitlabpls [u/b] myvariables`. 
    
    The value supports `${GIT_HEAD}` placeholder which will be automatically replaced by the app with git 
    HEAD branch of the working directory.

3. Run ```go install``` to build binary and install it to `$GOPATH/bin`

    _Note: you need to do this step every time the config changes because the config file is 
embedded within the binary._

4. Run `gitlabpls -v` to check if it's installed correctly

## Usage

```
$ gitlabpls help

USAGE:
   gitlabpls [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   url, u      generate URL
   browser, b  generate URL and open it in your browser

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```
