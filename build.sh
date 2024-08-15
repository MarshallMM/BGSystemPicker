# Step 1: Retrieve the Git description with tags
GIT_DESCRIBE=$(git describe --tags --dirty --always)

# Step 2: Retrieve the abbreviated commit hash
GIT_COMMIT_HASH=$(git rev-parse --short HEAD)

# Step 3: Combine them, ensuring the hash is always included
GIT_DESCRIPTION="${GIT_DESCRIBE}-${GIT_COMMIT_HASH}"

# Step 4: Build the Go binary, embedding the Git description
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-X 'main.gitDescription=${GIT_DESCRIPTION}'"
