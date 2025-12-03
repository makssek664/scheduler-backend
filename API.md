# API
## Format
The data format is JSON.

## Endpoints
### /auth
Authenticate a user, here meaining either login or create user, as of right now.
#### Data
* Name - nonnull name of the user.
* Returns:
    * ID - unique id of the user.
    * Name - name of the user.
    * ORM data
### /events/add/{user_id}
Add an event to a user of `{user_id}`.
#### Data
* Name - nonnull name of the event.
* Description - nullable descritpion of the event
* Color - nullable color tied to the event (frontend)
* Date - nonnull date component (JS compatible) indicating when it occurs.
* Returns:
    * ID - unique id of the event.
    * UserID - user id tied to the event.
    * Name - name of the event
    * Description - description of the event
    * Color - color of the event
    * Date - date of the event
    * ORM data.
### /events/get/{user_id}/all
Get all events related to user of `{user_id}`
#### Data
* Returns:
    * array:
        * ID - unique id of the event.
        * UserID - user id tied to the event.
        * Name - name of the event
        * Description - description of the event
        * Color - color of the event
        * Date - date of the event
        * ORM data. 
### /events/get/{user_id}/{event_id}
Get event of `{event_id}` related to user of `{user_id}`
#### Data
* Returns:
    * ID - unique id of the event.
    * UserID - user id tied to the event.
    * Name - name of the event
    * Description - description of the event
    * Color - color of the event
    * Date - date of the event
    * ORM data. 
### /events/set/{user_id}/{event_id}
Set event of `{event_id}` related to user of `{user_id}`
If a data member is null or, skippable, it is ignored.
#### Data
* Name - skipabble name of the event
* Description - skipabble description of the event
* Color - skippable color of the event
* Date - skippable date of the event
* Returns:
    * ok - always true.
### /events/rm/{user_id}/{event_id}
Remove event of `{event_id}` related to user of `{user_id}`
#### Data
* Returns:
    * ok - always true.
### Failure points
All endpoints have the two same failure points, InternalServerError and BadRequest.



