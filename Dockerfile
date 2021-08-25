FROM alpine:3
COPY github-to-terraform /bin/github-to-terraform
ENTRYPOINT [ "github-to-terraform" ]
