from google.cloud import storage

def get_storage_bucket(name):
    storage_client = storage.Client()

    return storage_client.bucket(name)