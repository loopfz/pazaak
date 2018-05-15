#!/bin/bash

#DECK1="+-2,+-2,+-2,+-6,+-3,+-3,+-3,+-4,+-4,+-5"
#DECK1="-3,-3,-3,-4,-4,-4,-2,-2,-6,-5"

if [ -z "$DECK1" ]; then
    DECK1="auto"
fi

if [ -z "$DECK2" ]; then
    DECK2="auto"
fi

if [ -z "$PLAYER1" ]; then
    PLAYER1="./testagent2/testagent2"
fi

if [ -z "$PLAYER2" ]; then
    PLAYER2="./testagent2/testagent2"
fi

if [ -z "$COUNT" ]; then
    COUNT=100
fi

if [ -z "$ROUNDS" ]; then
    ROUNDS=3
fi

# use this to rig the deck and give a certain value naturally to player 1
# this combined with "agentstand" implementations which just draw and stand
# can help you crunch natural drawing odds
if [ -z "$P1VALUE" ]; then
    P1VALUE=0
fi

SCORE1=0
SCORE2=0

for i in $(seq $COUNT)
do
    echo Running game $i...
    ./pazaakcli --quiet -player $PLAYER1 -player $PLAYER2 --round-limit $ROUNDS -p1-force-value $P1VALUE <<EOF
$DECK1
$DECK2
EOF
    if [ $? -eq 1 ]; then
        SCORE1=$(($SCORE1 + 1))
    else
        SCORE2=$(($SCORE2 + 1))
    fi
done

echo "score player 1 ($PLAYER1 [$DECK1]): $SCORE1"
echo "score player 2 ($PLAYER2 [$DECK2]): $SCORE2"
