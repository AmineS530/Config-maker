#!/bin/zsh

# ANSI color codes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
ORANGE='\033[1;38;5;100m'
RESET='\033[0m' # No Color

# Prompt for user's name and email
echo -e "${GREEN}Enter your fullname or login:${RESET} \c"
read NAME
echo -e "${GREEN}Enter your email:${RESET} \c"
read EMAIL

# Confirm the input
echo -e "\n${YELLOW}You have entered the following details:${RESET}"
echo -e "${GREEN}Name :${RESET} $NAME"
echo -e "${GREEN}Email:${RESET} $EMAIL"
echo -e "${YELLOW}Is this correct? (y/n):${RESET} \c"
read CONFIRMATION

# Check confirmation
if [[ $CONFIRMATION == "y" || $CONFIRMATION == "Y" ]]; then
  # Git configuration
  git config --global credential.helper store
  git config --global user.email "$EMAIL"
  git config --global user.name "$NAME"

  # Confirmation message
  echo -e "\n\033[1:0;32mGit has been configured with the following details:${RESET}"
  echo -e "${GREEN}Name :${RESET} $NAME"
  echo -e "${GREEN}Email:${RESET} $EMAIL"
else
  echo -e "\n${ORANGE}Aborted. No changes have been made.${RESET}"
fi

