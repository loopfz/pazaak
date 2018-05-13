#!/bin/bash

while :
do
    ./pazaakcli --quiet -player ./testagent/testagent2 -player ./testagent/testagent2 -stats pazaak_stats.json <<EOF
auto
auto
EOF
done
