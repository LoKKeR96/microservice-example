# Create a stage for building the application.

ARG GO_VERSION=1.23.5
# ARG BUILDPLATFORM=linux/amd64

# This image is made from alpine, this isn't the best linux image
# but it's small and helpful for quick development and running CI 
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS final_dev

WORKDIR /app

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage bind mounts to go.sum and go.mod to avoid having to copy them into
# the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=./src/go.mod,target=./go.mod \
    go mod download -x

# This is the architecture you're building for, which is passed in by the builder.
# Placing it here allows the previous steps to be cached across architectures.
# ARG TARGETARCH
RUN echo "I'm building for $TARGETPLATFORM"

# Copy the source code to the "usr" folder for quick rebuild
COPY ./src /usr/src

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage a bind mount to the current directory to avoid having to copy the
# source code into the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=./,source=./src \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server .

# Keep using the same docker image as base to enable development use cases
# like testing. 

# Expose the port that the application listens on.
EXPOSE 8080

# What the container should run when it is started.
ENTRYPOINT [ "/bin/server" ]
