FROM jetty

ADD /prediction-service-builder /usr/share/steam/.

WORKDIR /usr/share/steam/prediction-service-builder

RUN /bin/sh -c script.sh

#CMD ["bin/sh","script.sh"]
