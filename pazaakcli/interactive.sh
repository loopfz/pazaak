#!/bin/bash

cat - >/dev/tty

read -p "Move: " move </dev/tty >/dev/tty

echo $move
