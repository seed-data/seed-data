# Use the latest node lts version
FROM node:6.9.2

# Set the working directory & import the app
WORKDIR /app
COPY . /app

# Install the packages we need
RUN npm install
RUN npm install -g pushstate-server

# Build the server for production mode,
# this creates the /app/build folder with the compiled assets we need
RUN npm run build

# Expose port 80 for the server
EXPOSE 80

# Start a new instance of the server from /app/build on port 80:
# http://docker_ip:80
CMD pushstate-server build 80
