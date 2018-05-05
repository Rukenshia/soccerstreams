import os
import praw

def get_reddit_client():
    reddit = praw.Reddit(client_id=os.environ['REDDIT_CLIENT_ID'],
                        client_secret=os.environ['REDDIT_CLIENT_SECRET'],
                        user_agent=os.environ['REDDIT_USER_AGENT'],
                        username=os.environ['REDDIT_USERNAME'],
                        password=os.environ['REDDIT_PASSWORD'])

    if reddit.read_only:
        raise Exception('invalid credentials, read only mode will not be sufficient')
    return reddit