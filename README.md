# Twig
Twig is a `php` and `composer` CLI commands which will helps you proxy any commands for `php` and `composer` based on your project requirements. Twig will chose the right version for the PHP version when you run it inside of your projects, so you can have multiple PHP versions on your local machine without switching back and forth as you develop your application.

## List of commands
* `twig config init`

    To initialize the twig config directory


* `twig config set-default {version}`

    To set the current active php version


* `twig config --reload`

  To re-load the configuration file and sync to the newest changes. This command will invalidate caches and re-populate it with the newest content of the configuration file


* `twig php [...args]`

    To run the `php` command line with their arguments


* `twig composer [...args]`

    To run the `composer` command with their arguments

## Installation
Just download the pre-compiled binary from the [release page](https://github.com/dwadp/twig/releases) according to your operating system, and place it within your desired folder.

### Linux
Run the following command to download the executable binary:
```sh
curl -o {release_url}
```

Add the path to the binary on your `.bashrc` or `.zshrc` or any of your shell configuration file. Or if you prefer to use the standard binary location, you can put it in `/usr/local/bin` directory:
```sh
mv twig /usr/local/bin
```