import os

from prometheus_client import Histogram
from dotenv import load_dotenv
load_dotenv()

from reddit import get_reddit_client
from storage import get_storage_bucket
from comment import comment_agent
from post import post_agent

def main():
    metrics_event_diff = Histogram('praw_event_diff', 'PRAW Event diff', [], buckets=[20.0, 90.0])
    reddit = get_reddit_client()
    password = os.environ['RABBITMQ_PASSWORD']
    bucket = get_storage_bucket(os.environ['BUCKET_NAME'])

    if os.environ['AGENT_TYPE'] == 'post':
        post_agent(reddit, password, bucket, metrics_event_diff)
    elif os.environ['AGENT_TYPE'] == 'comment':
        comment_agent(reddit, password, bucket, metrics_event_diff)


if __name__ == "__main__":
    main()