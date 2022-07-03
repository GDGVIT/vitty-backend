FROM ubuntu:latest

WORKDIR /app

RUN apt-get update
RUN apt install -y python3-pip python3

ENV TZ=Asia/Kolkata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone


COPY . /app

RUN pip3 install fastapi uvicorn starlette
RUN apt install -y python3
RUN pip3 install -r requirements.txt
RUN pip3 install python-multipart

CMD ["uvicorn","main:app","--reload","--port","8000","--host","0.0.0.0"]
