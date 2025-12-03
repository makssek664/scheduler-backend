import requests
import json
import datetime
import unittest
class Mock(unittest.TestCase):
    def __init__(self, target: str, port: str):
        self.target = target
        self.port = port
        self.url = "http://" + target + ":" + port
    def __post_and_check__(self, endpoint, d):
        jsonb = d
        resp = requests.post(self.url + endpoint, json=jsonb)
        if resp.status_code != 201 and resp.status_code != 200:
            print(f"{resp.text} [{resp.status_code}]")
            self.fail(f"API CALL FAILED {resp.status_code}")
        return resp.json()
    def auth(self, username: str) -> int:
        data = {"Name": username}
        resp = self.__post_and_check__("/auth", data) 
        self.uid = resp["ID"]
        print(f"auth: {self.uid=}")
        return self.uid
    def addEvent(self, k):
        resp = self.__post_and_check__("/events/add/" + str(self.uid), k)
        print(f"addEvent: {resp['ID']=}")
        return resp["ID"]
    def getEvent(self, id):
        resp = self.__post_and_check__("/events/get/" + str(self.uid) + "/" + str(id), None)
        print(f"getEvent: {resp['ID']=} {resp['Name']=}")
        return resp
    def rmEvent(self, id):
        self.__post_and_check__("/events/rm/" + str(self.uid) + "/" + str(id), None)
    def getAllEvents(self):
        resp = self.__post_and_check__("/events/get/" + str(self.uid) + "/all", None)
        return resp
