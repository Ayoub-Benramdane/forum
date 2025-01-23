# Use Golang image with Alpine
FROM golang:1.22.2-alpine

# Install bash (if needed)
RUN apk add --no-cache bash gcc musl-dev

# Set working directory
WORKDIR /Forum

# Copy the application code
COPY . .

# Set the label for version
LABEL Version="1.0"

# Expose port for the application
EXPOSE 8404

RUN docker image prune -a


# Default command to run the application
CMD ["go", "run", "."]