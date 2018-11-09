Connect starts 4 routines [how they are stopped] :
    - heartbeat : periodically sends a heartbeat to the Gateway [error|stop]
    - listen : listens for payloads sent by the receiver or for a stop signal [error|stop]
        - receiver : receives message from the websocket connection and send them to the listen routine [error|websocket closure]
    - wait : monitors other routines, waiting for and error to occur to refresh the connection or for a stop signal [stop|multiple errors]


Consider using a special syntax for around/before/after parameters since they are mutually exclusive such as ~<> (maybe defaulting to after if none is specified).
Example : ~0123456789 for around 0123456789 or >0123456789 for after 0123456789.