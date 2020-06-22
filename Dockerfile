FROM golang:latest as prepare
 
WORKDIR /source
 
COPY go.mod .
COPY go.sum .
 
RUN go mod download
 
#
# STAGE 2: build
#
FROM prepare AS build
 
COPY . .
 
RUN CGO_ENABLED=0 GOOS=linux go build --tags dev -a -installsuffix cgo -o main .
 
#
# STAGE 3: run
#
FROM alpine:3.10 as run

COPY --from=build /source/main /main
COPY --from=build /source/init.sh /init.sh
 
RUN chmod +x /init.sh
 
ENTRYPOINT ["/init.sh"]
