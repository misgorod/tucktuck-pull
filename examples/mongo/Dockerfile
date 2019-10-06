FROM mongo
ADD mongo.conf /conf/mongo.conf
VOLUME [ "/data/db" ]
EXPOSE 27017
CMD mongod --config /conf/mongo.conf