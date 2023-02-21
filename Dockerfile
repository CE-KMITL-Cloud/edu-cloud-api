#########################################
#  Educational Cloud Platform Service   #
#    Computer Engineering, KMITL        #
#########################################

# Base Image
FROM golang:1.20.1-alpine

# Set Environment Variables
ENV PATH="/usr/local/go/bin:${PATH}"
ENV LANG='en_US.UTF-8' LC_ALL='en_US.UTF-8'
ENV TZ='Asia/Bangkok'

# Set working directory
WORKDIR /usr/src/app

# Coping source code to working directory
COPY . .

# Installs Go dependencies
RUN go mod download && go mod verify
RUN go mod tidy

# Buidling the Go app
RUN go build -o /usr/local/bin/app

EXPOSE 3001
CMD ["app"]