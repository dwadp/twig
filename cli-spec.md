# TODO

### Add cache mechanism
When a user run `twig php` for the first time, the app should get a matching PHP version from the configuration stored in `~/.twig/config.yaml` file, then it reads the `composer.json` or `.twig` file in the current directory, then it will match which PHP version to use according to the configuration files. After the first read, the PHP executable path should be cached for better performance so it shouldn't read the configuration file every time the `twig php` command runs. 

### List of commands
* `twig config --init`

    To initialize the twig config directory


* `twig config --set-default {version}`

    To set the current active php version


* `twig config --reload`

  To re-load the configuration file and sync to the newest changes


* `twig php [...args]`

    To run the php command line with their arguments


* `twig composer [...args]`

    To run the composer command with their arguments