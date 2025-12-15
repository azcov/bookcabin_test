FROM golang:1.25-alpine AS builder

RUN apk update && \
    apk upgrade && \
    apk add curl gcc g++ make git --no-cache \ 
    # && curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    # && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air \
    && mkdir -p /app 

WORKDIR /app
COPY    . .

EXPOSE 3000
CMD ["sh", "go","air"]