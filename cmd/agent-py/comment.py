from datetime import datetime
import json

from prometheus_client import Counter, start_http_server

from rabbitmq import get_rabbitmq_channel

# METRICS
PRAW_COMMENTS_INGESTED = Counter('praw_comments_ingested', 'PRAW Comments Ingested', ['parsing'])

def comment_agent(reddit, rabbitmq_password, bucket, metrics_event_diff):
    start_http_server(9000, addr='0.0.0.0')

    con = get_rabbitmq_channel(rabbitmq_password)
    con.queue_declare(queue='comments')
    con.close()
    
    # stream subreddit posts
    for comment in reddit.subreddit('soccerstreams').stream.comments():
        diff = datetime.utcnow() - datetime.utcfromtimestamp(comment.created_utc)

        blob = bucket.blob(comment.id)
        blob.upload_from_string(json.dumps({
            'id': comment.id,
            'author': comment.author.name,
            'author_flair_text': comment.author_flair_text,
            'body': comment.body,
            'created_utc': comment.created_utc,
            'link_id': comment.link_id,
            'name': comment.name,
            'parent_id': comment.parent_id,
            'ups': comment.ups,
        }))

        con = get_rabbitmq_channel(rabbitmq_password)
        con.basic_publish(exchange='',
                            routing_key='comments',
                            body=json.dumps({
                                "id": comment.id,
                                "created": comment.created_utc,
                            }))
        con.close()

        PRAW_COMMENTS_INGESTED.labels('skipped').inc()
        metrics_event_diff.observe(diff.seconds)
