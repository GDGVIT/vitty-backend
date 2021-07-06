FROM ubuntu:latest

WORKDIR /app

RUN apt-get update
RUN apt install -y python3-pip python3
RUN apt install -y libleptonica-dev

ENV TZ=Asia/Kolkata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt install -y pipenv tesseract-ocr
RUN apt install -y libtesseract-dev
RUN apt install -y python3-opencv

COPY . /app

RUN pip3 install Image Pillow
RUN pip3 install pytesseract
RUN pip3 install fastapi uvicorn opencv-python starlette
RUN pipenv install
RUN pip3 install python-multipart

CMD ["uvicorn","main:app","--reload","--port","8000","--host","0.0.0.0"]
