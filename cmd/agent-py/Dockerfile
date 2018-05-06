FROM python:2.7
ADD requirements.txt .

RUN pip install -r requirements.txt

ADD agent.py .
ADD reddit.py .
ADD comment.py .
ADD post.py .
ADD rabbitmq.py .
ADD storage.py .

ENV GOOGLE_APPLICATION_CREDENTIALS="/opt/soccerstreams/gcloud/gcloud-service-account.json"

CMD ["python", "./agent.py"]