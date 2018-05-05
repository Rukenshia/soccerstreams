FROM python:2.7
ADD requirements.txt .

RUN pip install -r requirements.txt

ADD reddit.py .
ADD comment.py .

CMD ["python", "./comment.py"]