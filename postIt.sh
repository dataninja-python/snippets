#!/bin/bash

base="http://localhost:4000"

echo "This script posts to the provided address."

read -p "What is the address following http://localhost:4000/? " ans

addr="$base/$ans"
echo ""

echo "The full address is: $addr"

echo ""
sleep 0.25
curl -i -X POST "$addr"
sleep 0.25

echo ""

sleep 0.50

exit 0
