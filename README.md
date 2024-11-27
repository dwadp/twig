# Twig
> [!WARNING]  
> This is still an under development project, although there is finished project that you can clone and build yourself by checking the `main` branch, but i think it can be improved more. And in the end, when the project is consider to be stable, it will be released as version 2.

Twig is a `php` and `composer` CLI commands which will helps you proxy any commands for `php` and `composer` based on your project requirements. Twig will chose the right version for the PHP version when you run it inside of your projects, so you can have multiple PHP versions on your local machine without switching back and forth as you develop your application.

## Usage
### Initialize configuration file
The first thing you want to do after installing `twig` is to initialize the configuration file. You can do this by running the following command:
```sh
twig config init
```
Tipically the configuration file will be placed under your `$HOME` directory for example: `$HOME/.twig/twig.yaml`

### Set the PHP versions
You can set the default PHP version to be used when you're running the `php` command outside of your project directory.
```sh
twig config set-default {version}
```

### Running `php` command
This is the core function of twig, which is to run the `php` command inside of your project directory
```sh
twig php [args...]
```
By default, `twig` will read your `composer.json` file and will search for the `require` and then look for the `php` key to determine the PHP version to be used. If the PHP version is not found, it will use the default PHP version that you've set before.

### Running `composer` command
Another core function of twig is to run the `composer` command inside of your project directory
```sh
twig composer [args...]
```
The `composer` command will be proxied by `twig` and will use the PHP version that is set on your `composer.json` file. If the PHP version is not found, it will use the default PHP version that you've set before.


> [!TIP]
> To get more information about the commands, you can run `twig --help` or `twig {command} --help`

## Installation [WIP]
Just download the pre-compiled binary from the [release page](https://github.com/dwadp/twig/releases) according to your operating system, and place it within your desired folder.

### Pre-Built Binary [WIP]
#### Linux
Run the following command to download the executable binary:
```sh
curl -o {release_url}
```

Add the path to the binary on your `.bashrc` or `.zshrc` or any of your shell configuration file. Or if you prefer to use the standard binary location, you can put it in `/usr/local/bin` directory:
```sh
mv twig /usr/local/bin
```

### Build from Source
Although the pre-built binary is available, you can also build the binary from the source code. To do this, you need to have `go` installed on your machine at least version `v1.23.3`.

#### Clone repository
Run the following command to clone the repository:
```sh
git clone https://github.com/dwadp/twig.git
```

#### Build
Run the following command to build the binary:
```sh
VERSION=$(git rev-parse --short HEAD) CGO_ENABLED=0 && go build -ldflags "-s -w -X 'github.com/dwadp/twig/cmd/twig.Version=2.0.0-$VERSION'" -o twig main.go
```
Obviously you can tweak the flags however you want.

#### Place your binary
You can place the `twig` binary anywhere in your system, for example in `/usr/local/bin` directory:
```sh
mv twig /usr/local/bin
```

#### Verify the installation
To verify the installation, you can run the following command:
```sh
twig --version
```

### Configure aliases
To make it easier to run the `php` and `composer` using twig, you can add the following aliases to your `.bashrc` or `.zshrc` or any of your shell configuration file:
```sh
alias php='twig php'
alias composer='twig composer'
```