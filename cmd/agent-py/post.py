from datetime import datetime
import json

from prometheus_client import Counter, start_http_server

from rabbitmq import get_rabbitmq_channel


# METRICS
PRAW_POSTS_INGESTED = Counter('praw_posts_ingested', 'PRAW Posts Ingested', ['parsing'])

def post_agent(reddit, rabbitmq_password, bucket, metrics_event_diff):    
    start_http_server(9000, addr='0.0.0.0')

    con = get_rabbitmq_channel(rabbitmq_password)
    con.queue_declare(queue='posts')
    con.close()

    # stream subreddit posts
    for post in reddit.subreddit('soccerstreams').stream.submissions():
        diff = datetime.utcnow() - datetime.utcfromtimestamp(post.created_utc)

        con = get_rabbitmq_channel(rabbitmq_password)
        con.basic_publish(exchange='',
                            routing_key='posts',
                            body=json.dumps({
                                "id": post.id,
                                "created": post.created_utc,
                            }))
        con.close()

        PRAW_POSTS_INGESTED.labels('skipped').inc()
        metrics_event_diff.observe(diff.seconds)
        
