#!/usr/bin/env python3
import zmq
from datetime import datetime as dt
ctxt = zmq.Context()
pusher = ctxt.socket(zmq.PUSH)
pusher.set_hwm(1)
pusher.connect('tcp://127.0.0.1:2048')

t0 = tstamp = dt.now()
nr_sent = 0
while True:
    now = dt.now()
    pusher.send_unicode(str(now))
    nr_sent += 1
    if (now - tstamp).total_seconds() > 1:
        rate = nr_sent / ((now - t0).total_seconds())
        print('%d: %.0f msgs/sec' % (nr_sent, rate))
        tstamp = now 

