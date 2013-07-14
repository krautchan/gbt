gbt
===

IRC bot written in GO

#### Compilation & Installation [![Build Status](https://travis-ci.org/krautchan/gbt.png?branch=master)](https://travis-ci.org/krautchan/gbt)

    $ go get github.com/krautchan/gbt
    $ go install github.com/krautchan/gbt

#### Configuration
Run gbt once and it will create the directory `$HOME/.config/gbt`. If you run gbt as root(DON'T) it will use `/etc/gbt`. In that directory you will find the main config file `config.conf` and a subdirectory called `dev-urandom`.

In `config.conf` you define the servers you want your bot to connect to. For each server entry gbt will create a new subdirectory with all module configuration files.

Example:

    {"config":
        [
            {"name":"dev-urandom","address":"dev-urandom.eu","port":"6667"},
            {"name":"example","address":"example.com", "port":"6667"}
        ]
    }

When you added a new server you should restart gbt and stop it again. Go to the admin.conf file in the subdirectory of the added server and change the admin password to something more secure. Go through the other config files and change the setting to your liking.

#### Running
It is recommended to start gbt in a `screen` or `tmux` session. To get a list of all available commands type &help.
