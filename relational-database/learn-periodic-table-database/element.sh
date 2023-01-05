#!/bin/bash

if [[ -z $1 ]]
then
  echo "Please provide an element as an argument."
  exit 0
fi

PSQL="psql -X --username=freecodecamp --dbname=periodic_table --tuples-only -c"

ELEMENTS=$($PSQL "SELECT atomic_number, symbol, name FROM elements")
while read ATOMIC_NUMBER BAR SYMBOL BAR NAME; do
  if [[ $1 -eq $ATOMIC_NUMBER ]] || [[ $1 -eq $SYMBOL ]] || [[ $1 -eq $NAME ]]
  then
      ELEMENT_PROPERTY=$($PSQL "SELECT atomic_mass, melting_point_celsius, boiling_point_celsius, type FROM properties INNER JOIN types using (type_id) WHERE atomic_number = $ATOMIC_NUMBER")
      while read ATOMIC_MASS BAR MELTING_POINT_CELSIUS BAR BOILING_POINT_CELSIUS BAR TYPE; do
        echo "The element with atomic number $ATOMIC_NUMBER is $NAME ($SYMBOL). It's a $TYPE, with a mass of $ATOMIC_MASS amu. $NAME has a melting point of $MELTING_POINT_CELSIUS celsius and a boiling point of $BOILING_POINT_CELSIUS celsius."
      done < <(echo -e  "$ELEMENT_PROPERTY")
      exit 0
  fi
done < <(echo -e  "$ELEMENTS")

echo "I could not find that element in the database."
