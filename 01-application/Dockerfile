FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go build -o server .
ENV AUTHOR="John Doe"
EXPOSE 8000
CMD ["./server"]