FROM node:14

# install simple http server for serving static content
RUN npm install -g http-server

# Set the working directory
WORKDIR /visualizer-front

# Clone the repository
RUN git clone https://github.com/Javtor/challenge-zinc.git .

RUN git fetch

RUN git checkout docker

# Go to the appropiate folder
WORKDIR ./visualizer-front/challenge-zinc-visualizer

# Install the app's dependencies
RUN npm install

# Build the app for production
RUN npm run build

EXPOSE 8080

# Run the binary
CMD [ "http-server", "dist" ]
