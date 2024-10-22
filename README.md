# Public Ubuntu Configuration For Campus
## Informations:
 - This script was made to make setting up your session after changing it/resetting it easier/less painful
 - Feel free to fork and adjust it to your needs :)
      ### features:
      + Makes ZSH as your default terminal, can be undone and re done with `make bash`(to revert to bash as default) or `make zsh` afterwards 
      + Automatically adds french keyboard layout
      + changes logout time automatically
      + Changes terminal settings and defaults using zsh for better experience
      + Prompts to change wallpaper upon using, you can choose your own, or select one of the pre installed ones from the folder
      + Changes the theme based on preference
      + Sets up github/gitea for first time use and to remember your password whenever you use `git push/gp`
      + Changes default terminal font for compatibility with terminal theme
      + Set up docker easily with
        - `make docker`

## Requirements:
- zsh terminal
  
## How to use:
- Download current [release](https://github.com/AmineS530/Config-maker/releases) or `Makefile` from the [repository](https://github.com/AmineS530/Config-maker/blob/main/Makefile)
  `If downloading the Makefile from the repository make sure to remove the ending '.txt'`
- Run it using
```zsh
make
```

## Recent:
  - Added make [prompt] for script that can be reused
      + `make git` : redo github/gitea settings
      + `make theme` : change theme
      + `make background` : change background
  - Cleaned Makefile rules
  - Now upon finishing all instances of gnome-terminal close
  - Set logout to 1hour

## FAQ
- ### I didn't enter my name or email correctly during first time use:
    run the following command in same dir as the Makefile:

    ```zsh
    make git
    ```

### ToDo:
- More wallpaper options.
- Set-up a variable file that can be customized for your own settings after forking
- adding some 'funsies' such as [sl](https://github.com/thekakester/sl)
- Adding custom themes aside from the default ones
- adding VIM customization settings

## Feedback:
- If you have any ideas to improve this or have any sort of feedback, feel free to add and DM on Discord `a.sadik`.

# ***Happy coding! :)***
