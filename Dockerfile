# use an existing docker image as a base
FROM golang:alpine as builtapp
#Installing additional tools for the container
WORKDIR '/go/src/app'
COPY . .
RUN apk add --no-cache git mercurial \
   && go get -d -v github.com/tools/godep \
   && apk del git mercurial
RUN go install
#Commands, run at container startup
CMD ["app"]


#CMD COMMANDS------------------------------------------

#FROM alpine as release
#COPY --from=builtapp /out /out
#WORKDIR '/go/src/app'
#CMD ["app"]

# commands: docker run -p 8080:8080 -v "//c/Users/Public/Weather_app_development:/go/src/app" ae75d2ef6e30