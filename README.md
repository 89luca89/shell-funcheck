# Shell funchek

Having fun, checking shell function comments

![Shell Funcheck Logo](./assets/funcheck-logo.png)

## Idea

To improve shell script maintainability, having some comments will surely help.

Shell script don't have an easy way to keep track of global variables, env variables and arguments,
this project aims to "lint" functions and force the developer to document them above the said function.

This is particularly useful for [Distrobox](https://github.com/89luca89/distrobox)

## Compile

```sh
CGO_ENABLED=0 go build -mod vendor -ldflags="-s -w"
```

Will generate a statically compiled binary

## Usage

```log
Check shell functions for undocumented stuff

Usage:
  shell-funcheck [command]

Available Commands:
  check       Check for undocumented arguments and global variables inside functions
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -h, --help      help for shell-funcheck
  -v, --version   version for shell-funcheck

Use "shell-funcheck [command] --help" for more information about a command.
```

### Examples

Given the following script:

```sh
#!/bin/sh

global_var="hello"

foo() {
  input=$1
  echo $global_var $input from $USER
}

foo "world"
```

running `shell-funcheck check test.sh`

Will report us:

```log
test.sh:5:1: error: foo - foo: function foo should be documented
test.sh:6:3: error: foo - input: argument is not documented
2024/05/01 18:15:31 lint error
```

Which already gives us an hint **document the function**

### How to document the function

Expected documentation format is:

```sh
# name_of_funtion + explanation
# Arguments:
#   name_of_argument: explanation of argument
# Expected env variables:
#   name_of_env_variable: explanation
# Expected global variables:
#   name_of_global_variable: explanation
# Outputs:
#   explanation of the expected outputs
```

If we want to fix our previous example, we should have something like:

```sh
#!/bin/sh

global_var="hello"

# foo will print our nice hello world
# Arguments:
#   input: string world to salute
# Expected env variables:
#   USER: string of your username
# Expected global variables:
#   global_var: string global var of this script
# Outputs:
#   a nice salute
foo() {
  input=$1
  echo $global_var $input from $USER
}

foo "world"
```

## Credits

Many thanks to [mvdan](https://github.com/mvdan/) for their fantastic library [https://github.com/mvdan/sh](https://github.com/mvdan/sh) that
made this project feasable
