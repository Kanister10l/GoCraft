FROM openjdk:8-jre-alpine3.9
RUN mkdir /mc-root
RUN mkdir /mc
WORKDIR /mc-root
RUN wget -O mc.jar "https://cdn.getbukkit.org/craftbukkit/craftbukkit-1.14.4-R0.1-SNAPSHOT.jar"
COPY eula.txt ./eula.txt
COPY server.properties ./server.properties
COPY startscript.sh ./startscript.sh
RUN chmod 755 /mc-root/startscript.sh
CMD ["/mc-root/startscript.sh"]