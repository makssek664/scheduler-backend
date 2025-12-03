#!/bin/bash
. scripts/ticktick.sh
. ./mocks/add_event.sh

EVENTDATA2=$(curl -X POST localhost:8080/events/get/1/all)
echo $EVENTDATA2
