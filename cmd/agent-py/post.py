from reddit import get_reddit_client
from datetime import datetime

from prometheus_client import Histogram, Counter, start_http_server
from dotenv import load_dotenv
load_dotenv()

start_http_server(9000, addr='0.0.0.0')

# METRICS
PRAW_EVENT_DIFF = Histogram('praw_event_diff', 'PRAW Event diff', [], buckets=[20.0, 90.0])
PRAW_POSTS_INGESTED = Counter('praw_posts_ingested', 'PRAW Posts Ingested', ['parsing'])

def main():
    reddit = get_reddit_client()
    
    # stream subreddit posts
    for post in reddit.subreddit('soccerstreams').stream.submissions():
        diff = datetime.utcnow() - datetime.utcfromtimestamp(post.created_utc)

        PRAW_POSTS_INGESTED.labels('skipped').inc()
        PRAW_EVENT_DIFF.observe(diff.seconds)


if __name__ == "__main__":
    main()