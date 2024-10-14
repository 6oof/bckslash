#!/bin/bash

echo "Starting an action..."
# This script handles various actions. The general scaffolding of this file should not be changed; only the actions within the if statements should be added. 
[ -z "$1" ] && { echo "No action provided!"; echo "Usage: ./script.sh <action>"; echo "Available actions: deploy, status, boot, destroy"; exit 1; } 

if [ "$1" == "deploy" ]; then
    echo "Deploying the application..."
    # Add deploy logic here
    sudo ls
elif [ "$1" == "status" ]; then
    echo "Checking project status..."
    # Add status check logic here

elif [ "$1" == "boot" ]; then
    echo "Booting the system..."
    # Add boot logic here

elif [ "$1" == "destroy" ]; then
    echo "Destroying resources..."
    # Add destroy logic here

fi
