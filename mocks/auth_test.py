#!/bin/python3
import qxm_mock

m = qxm_mock.Mock("localhost", "8080")
m.auth("TestUser")
