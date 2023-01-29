from cassandra.cluster import Cluster
from . import configs


class CassandraColumn:

    def __init__(self, config: configs.CassandraConfig):

        self.config = config

        cluster = Cluster([config.cassandra_listen_address])

        self.session = cluster.connect(config.cassandra_keyspace)
