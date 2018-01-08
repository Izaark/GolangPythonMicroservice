FROM izaark/go-gingonic
COPY . /go/src/github.com/tasks/GolangPythonMicroservice/Go
RUN cp /go/src/github.com/tasks/GolangPythonMicroservice/Go/*.env /go
RUN go install github.com/tasks/GolangPythonMicroservice/Go/
RUN rm -rf /go/src/github.com/tasks/GolangPythonMicroservice/Go/
CMD Go
EXPOSE 8080
