FROM golang:1.15-alpine
ENV PACKAGES make git curl build-base gcc musl-dev linux-headers
WORKDIR /app
RUN apk add --update $PACKAGES
COPY ./bin/mytokenapp  ./
#RUN chmod 0777  /usr/local/bin/mytokenapp
#RUN apk add --update ca-certificates
#RUN apk add --no-cache bash
EXPOSE 1317 26656 26657
CMD ["mytokenapp"]





