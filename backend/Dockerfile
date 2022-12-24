FROM golang:1.20rc1-alpine3.17
# working directory
WORKDIR /app
# copy project and download modules
COPY . ./
RUN go mod download
# build
RUN go build -o /backend

CMD [ "/backend"]