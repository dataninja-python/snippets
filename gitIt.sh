#!/bin/bash

git add .

read -p "What is the git message? " ans
echo "You entered: $ans"

git commit -am "$ans"

git push origin main
