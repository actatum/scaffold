FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN echo "[url \"git@github.com:\"]\n\tinsteadOf = https://github.com/" >> /root/.gitconfig

# add credentials on build
ARG SSH_PRIVATE_KEY
RUN mkdir /root/.ssh/
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa

# make sure your domain is accepted
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN chmod 400 /root/.ssh/id_rsa

RUN cd cmd/server && go build -o main

# Remove SSH keys
RUN rm -rf /root/.ssh/

EXPOSE 8080

CMD ["./cmd/server/main"]