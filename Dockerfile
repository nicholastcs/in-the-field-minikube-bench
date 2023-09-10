FROM golang as builder
WORKDIR /workspace
COPY . .

RUN go get .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
WORKDIR /bin/

COPY --from=builder /workspace/app .
CMD [ "./app" ]