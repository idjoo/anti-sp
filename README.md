# anti-sp
Program created for binusian to minimalize the chance to fail classes because of
absence, written purely in go.

## Installation
    go install github.com/cocatrip/anti-sp@latest

## Usage
When first running anti-sp will ask for a username and password:

    Username: <Your Bimay Username (without @binus.ac.id)>
    Password: <Your Bimay Password>

Example:

    Username: rayhan.noersandi
    Password: mySecurePass!@#

### Credential
Your username and password will be stored locally in your computer inside the UserConfigDir. UserConfigDir returns the default root directory to use for user-specific configuration data.
Users should create their own application-specific subdirectory within this one and use that.

On Unix systems, it returns $XDG_CONFIG_HOME as specified by [freedesktop](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) if non-empty, else $HOME/.config. On Darwin, it returns $HOME/Library/Application Support. On Windows, it returns %AppData%. On Plan 9, it returns $home/lib.


