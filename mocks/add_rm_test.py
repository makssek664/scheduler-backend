import qxm_mock
from datetime import datetime
import json

m = qxm_mock.Mock("localhost", "8080")
m.auth("TestUser")
preadd_len = len(m.getAllEvents())
id = m.addEvent({"Name": "abc", "Date": datetime.now().strftime('%Y-%m-%dT%H:%M:%SZ')})
ev = m.getEvent(id)
print(ev)
m.rmEvent(id)
if len(m.getAllEvents()) != preadd_len:
    print("Backend did not remove event!")
