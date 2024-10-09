FROM golang:1.22.2-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the rest of the application code
COPY . .

RUN mkdir -p /app/var/logs && touch log.txt && chmod -R 755 /app/var/logs 

RUN ls -lrt 

RUN CGO_ENABLED=0 go build -o pismoapp .

RUN chmod +x /app/pismoapp

#build image 
FROM alpine:latest

RUN mkdir /app
#copy from builder to 
COPY --from=builder /app/pismoapp /app   

CMD [ "/app/pismoapp" ]

