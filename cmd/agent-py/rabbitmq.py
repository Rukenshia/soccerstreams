import pika
import os

def get_rabbitmq_channel(password):
    connection = pika.BlockingConnection(pika.ConnectionParameters('esteemed-frog-rabbitmq.default.svc.cluster.local', heartbeat=0, credentials=pika.credentials.PlainCredentials('user', password)))
    return connection.channel()