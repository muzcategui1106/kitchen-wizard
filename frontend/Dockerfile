FROM node:18.5.0-slim as builder

# set working directory
WORKDIR /app

# pre-warm dependencies
RUN npm install react-scripts@3.4.1 -g --silent

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# install app dependencies
COPY package.json ./
COPY package-lock.json ./
RUN npm install --silent

# add app
COPY . ./

# start app
CMD ["npm", "start", "--host", "0.0.0.0"]