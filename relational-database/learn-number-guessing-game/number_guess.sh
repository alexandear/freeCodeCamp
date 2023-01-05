#!/bin/bash

PSQL="psql --username=freecodecamp --dbname=number_guess -t --no-align -c"

MIN_NUMBER=1
MAX_NUMBER=1000

SECRET_NUMBER=$(( RANDOM % $MAX_NUMBER + $MIN_NUMBER ))

echo "Enter your username:"
read USERNAME

PLAYER_ID=$($PSQL "SELECT player_id FROM players WHERE username = '$USERNAME'")
if [[ -z $PLAYER_ID ]]
then
  INSERT_PLAYER_RESULT=$($PSQL "INSERT INTO players(username) VALUES('$USERNAME')")
  PLAYER_ID=$($PSQL "SELECT player_id FROM players WHERE username = '$USERNAME'")
  echo "Welcome, $USERNAME! It looks like this is your first time here."
else
  GAMES_PLAYED=$($PSQL "SELECT COUNT(*) FROM games WHERE player_id = $PLAYER_ID")
  BEST_GUESS=$($PSQL "SELECT MIN(guess) FROM games WHERE player_id = $PLAYER_ID")
  if [[ -z $BEST_GUESS ]]
  then
    BEST_GUESS=0
  fi

  echo "Welcome back, $USERNAME! You have played $GAMES_PLAYED games, and your best game took $BEST_GUESS guesses."
fi

echo "Guess the secret number between $MIN_NUMBER and $MAX_NUMBER:"

NUMBER_OF_GUESSES=1
while : ; do
    read GUESS_NUMBER

    if [[ ! $GUESS_NUMBER =~ ^[0-9]+$ ]]
    then
      echo "That is not an integer, guess again:"
      continue
    fi

    if [[ $GUESS_NUMBER -lt $SECRET_NUMBER ]]
    then
      echo "It's lower than that, guess again:"
    elif [[ $GUESS_NUMBER -gt $SECRET_NUMBER ]]
    then
      echo "It's higher than that, guess again:"
    else
      break
    fi

    ((NUMBER_OF_GUESSES++))
done

INSERT_GAME_RESULT=$($PSQL "INSERT INTO games(player_id, guess) VALUES('$PLAYER_ID', $NUMBER_OF_GUESSES)")

echo "You guessed it in $NUMBER_OF_GUESSES tries. The secret number was $SECRET_NUMBER. Nice job!"
