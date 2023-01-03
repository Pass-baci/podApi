FROM alpine
COPY podApi /podApi
WORKDIR "/"
CMD ["./podApi"]