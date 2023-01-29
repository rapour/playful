import numpy as np
from hmmlearn import hmm

from tools.cassandra import CassandraColumn
from tools.configs import CassandraConfig


def learn():

    
    top = np.array([])

    csConfig = CassandraConfig()
    cs = CassandraColumn(csConfig)

    rows = cs.session.execute("SELECT * FROM locations")
    for row in rows:
        top = np.append(top, [int(row.altitude)])

    
    top = top.astype(np.int64)
    if not np.issubdtype(top.dtype, np.integer):
        raise ValueError(f"Data type is not correct: {top.dtype}")

    if top.min() < 0:
        raise ValueError(f"Minimum is below zero: {top.min()}")

    
    top = top.reshape(-1, 1)
    model = hmm.MultinomialHMM(n_components=200, n_trials=top.sum(axis=1))

    print(model.fit(top))
    