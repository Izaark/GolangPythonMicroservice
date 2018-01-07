FROM izaark/python-flask
COPY . /Python
WORKDIR /Python


RUN pip3 install -r requirements.txt
RUN pip3 install flask
ENTRYPOINT ["python"]
CMD ["main.py"]
EXPOSE 4000