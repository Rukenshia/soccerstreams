from reddit import get_reddit_client
from datetime import datetime

from prometheus_client import Histogram, Counter, start_http_server
from dotenv import load_dotenv

import os
print(os.environ)


# METRICS
PRAW_EVENT_DIFF = Histogram('praw_event_diff', 'PRAW Event diff', [], buckets=[20.0, 90.0])
PRAW_COMMENTS_INGESTED = Counter('praw_comments_ingested', 'PRAW Comments Ingested', ['parsing'])

def main():
    reddit = get_reddit_client()

    start_http_server(9000, addr='0.0.0.0')
    
    # stream subreddit posts
    for comment in reddit.subreddit('soccerstreams').stream.comments():
        diff = datetime.utcnow() - datetime.utcfromtimestamp(comment.created_utc)

        PRAW_COMMENTS_INGESTED.labels('skipped').inc()
        PRAW_EVENT_DIFF.observe(diff.seconds)


if __name__ == "__main__":
    main()