#!/bin/sh

cd /mc

if [ "$(ls -A /mc)" ]; then
        echo "Running previous instance of server"
        java -jar /mc/mc.jar
else
        echo "Creating new instance of server"
        cp /mc-root/server.properties /mc/server.properties
        cp /mc-root/eula.txt /mc/eula.txt
        cp /mc-root/mc.jar /mc/mc.jar
        java -jar /mc/mc.jar
fi