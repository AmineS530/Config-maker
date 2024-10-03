# Public Ubuntu Configuration For Campus
### ToDo:
- make [prompt] `as in make git to execute the git setup script without navigating to the script manually`.
- theme script can be confusing as it displays full name of the theme which can be confusing, also make it to get name of the themes depending on the mashine rather than what's predefined.
- More wallpaper options.
- set-up a variable file that can be customized for your own settings after forking
## Informations :
   - This script was made to make setting up your session after changing it/resetting it easier/less painful
   - feel free to fork and adjust it to your needs :)
        ####    features
        + Changes terminal settings and defaults using zsh for better experience
        + Prompts to change wallpaper upon using, you can choose your own, or select one of the pre installed ones from the folder
        + Changes the theme based on preference
        + Sets up github/gitea for first time use and to remember your password whenever you use `git push` or `gp`
        + changes default terminal font for compatibility with terminal theme
        + sets up docker easily with `make docker`

## Requirements :
- zsh terminal + Oh My Zsh installed.
  
## Usage :
- Download current [release](https://github.com/AmineS530/Config-maker/releases) or `Makefile` from the [repository](https://github.com/AmineS530/Config-maker/blob/main/Makefile)
  `If downloading the makefile from the repository make sure to remove the ending '.txt'`
- Run it using
```zsh
make
```

## Tips :
- In case of encountering any error or failing anything, you can go to `~/home/Zone01_Desk_cfg/` and reuse the script you need to reapply using
```zsh
zsh [stript_name].sh
```

## FAQ

#### I didn't enter my name or email correctly during first time use

go to `~/Zone01_Desk_cfg` and execute the following:
```zsh
zsh git_setup.sh
```

## Feedback :
- If you have any ideas to improve this or have any sort of feedback, feel free to DM me on Discord (username: a.sadik).

## **Happy coding!** :)
