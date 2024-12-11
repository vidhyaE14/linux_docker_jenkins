FROM golang:alpine3.21
LABEL Name="Vidhya"
RUN apk add --no-cache git
RUN git clone https://github.com/vidhyaE14/goLang_project.git
#RUN cd goLang_project
WORKDIR /go/goLang_project
ENV DB_USERNAME admin
ENV DB_PASSWORD Admin_1234
ENV RDS_ENDPOINT inventorysystem.cjuei8gaazd2.us-east-2.rds.amazonaws.com
ENV RDS_NAME inventorysystem
ENTRYPOINT ["go","run","main.go"]
EXPOSE 8081