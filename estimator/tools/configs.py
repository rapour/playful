from pydantic import BaseSettings



class CassandraConfig(BaseSettings):
    cassandra_listen_address: str
    cassandra_keyspace : str

    class Config:
        case_sensitive = False
        allow_mutation = False
        env_prefix = ''
