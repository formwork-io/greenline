#!/usr/bin/env python3
import zmq
from datetime import datetime as dt
ctxt = zmq.Context()
puller = ctxt.socket(zmq.PULL)
puller.set_hwm(1)
puller.bind('tcp://127.0.0.1:8402')

t0 = tstamp = dt.now()
nr_recvd = 0
while True:
    puller.recv_unicode()
    nr_recvd += 1
    now = dt.now()
    if (now - tstamp).total_seconds() > 1:
        rate = nr_recvd / ((now - t0).total_seconds())
        print('%d: %0.0f msgs/sec' % (nr_recvd, rate))
        tstamp = now 

