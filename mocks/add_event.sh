#!/bin/bash
. scripts/ticktick.sh

USERDATA=$(./mocks/auth.sh)
tickParse "$USERDATA"
EVENTDATA=$(
    curl -X POST -d '{"name": "something has to give", "date": "2001-09-11T11:11:11.000Z"}' localhost:8080/events/add/``ID``)
    echo $EVENTDATA
