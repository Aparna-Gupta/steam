FROM jetty

ADD /prediction-service-builder /usr/share/steam/.

WORKDIR /usr/share/steam/prediction-service-builder

#RUN script.sh

CMD ["bin/sh","script.sh"]
