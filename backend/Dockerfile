FROM golang:1.19.4-bullseye
# working directory
WORKDIR /app
# copy project and download modules
COPY . ./
RUN go mod download
# build
RUN go build -o /backend

CMD [ "/backend"]