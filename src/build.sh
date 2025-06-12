# Step 1: Retrieve the Git description with tags
GIT_DESCRIBE=$(git describe --tags --dirty --always)

# Step 2: Build the Go binary, embedding the Git description
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-X 'main.gitDescription=${GIT_DESCRIBE}'"
sleep 2