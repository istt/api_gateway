# SMS Platform Traffic Simulator User Manual

### Introduction

Idea of the traffic simulator application is, that we can connect this application via SMPP to our SMS Platform on input and output side. This means, that it allows us to simulate to send SMS like a customer, and it will simulate to receive SMS and to send DLR responses like a supplier.

The traffic simulator will then send SMS patterns on input connections and receives the SMS back over the output connections. SMS received on output connections will be replied with certain DLR patterns.

### Terms and Abbreviation

- ESME Account: SystemID, Password, IP, Port to establish a SMPP session
  to SMSC.
- ESME Session: a persistence connection between ESME towards SMSC. One
  account can have multiple sessions. By persist, the ESME tcp socket
  won't change during connection
- SMPP Test Case: Combination of SMPP Account, Num Sessions, SMPP
  Parameters, and Receiver List
- Receiver List: list of receivers, normally MSISDN
- MT Submit: ESME → Horisen SMSGW → SMSC (Supplier) (Acknowledgement
  Accept ratio)
- MO Deliver: SMSC → Horisen SMSGW → ESME based on short code
- DLR Deliver Receipt: ESME → Horisne SMSGW → SMSC (delay max 1 day) →
  Horisen SMSGW → ESME (random delivery status)

### Overview

In a later phase, also part of simulator is to give answer on MNP
queries, to be able to simulate MNP answers as well, but again, this in
2nd phase of the project.

#### Who will use this Traffic Simulator, what are the business cases?

1.  Feeding "real-time" traffic data into our Demo Platform, to have all
    the time "good" and relevant data to show during sales and training
    process on our Demo Platform.  
    Here we need a repeating cronjob, to have 24/7/365 data available to
    simulate real traffic.

2.  Sales people can trigger a sending task to show to clients, how fast
    our platform works and how the monitoring tools works in real-time.
    Sales will trigger "Traffic Peaks" which are visible very fast on
    the statistics.

3.  Routing Test: based on SMS sent by Traffic Simulator app, the user
    can test his routing settings made in platform and see if his
    setting are working accordingly.

4.  Performing stress tests that last long period of time with amounts
    of traffic that vary over the time.

#### Features Sending Traffic / Customer Simulator

1.  Supported Protocols for connecting:

    1.  SMPP

    2.  HTTP

2.  Send SMS one by one over simple Sending Mask with all parameters
    defined.

3.  Send SMS via pattern imported via CSV/Excel

The sending pattern is taken from a file, where the simulator just
repeats this sending. In this file all relevant data is provided, to be
able that simulator can send the traffic (sender, destination, text,
account, timestamp, etc.).

A sending pattern can contain messages for few seconds up to a complete
month

Features:

- Multiple sending patterns can be imported into the simulator

- Many sending patterns can run in parallel (overlapping)

  - This means that at the same time many patterns can be run, like
    manual sending patterns, repeating patterns, etc.

- Cronjob with repeating sending of certain pattern (also with
  start/stop/suspend actions)

  - Repeating a pattern will send constantly traffic, without any
    break

  - Example: if pattern lasts 24h, it will be repeated continuously
    every 24h

- START/STOP mechanism to start/stop sending certain pattern

- Pattern variation in a daily % range to change the values in logs
  and stats (described in details below)

<!-- end list -->

4.  Sending traffic pattern with defined Throughput (stress test)

#### Features Receiving Traffic / Suplier Simulator

1.  Supported Protocols for connecting to Traffic Simulator:

    1.  SMPP

    2.  HTTP (later phase)

2.  The Traffic Simulator needs to "play" the role of a supplier and
    should respond accordingly for the traffic received.  
    Following parameters can be defined:

    3.  DLR traffic pattern based on imported traffic pattern

    4.  Generated responses based on parameters:

        1.  DLR Ratio Range (e.g. 84-98%) per day

        2.  DLR Reply Range (e.g. 0.5sec – 16sec and 50% @ 4sec) per
            message

        3.  Define DLR Error Codes to send back

        4.  Reject Ratio Range (e.g. 0.03-2.4%) per day

        5.  Define Reject Error Codes to send back

### Customer simulator details

Customer simulator requires as input:

1.  List of SMPP and HTTP accounts with credentials (excel CSV file)
2.  Connection JSON file - SMPP address/port and HTTP URL for connecting and sending traffic
3.  Traffic template (excel CSV file)

#### List of accounts with credentials

This is Excel CSV file that contains:

1.  Column caption row (one)
2.  Account rows

Each account row has following columns:

- Protocol ("smpp" or "http")
- Username (or system_id for SMPP)
- Password
- System type (for SMPP)
- Number of binds (for SMPP) – or Number of possible parallel
  submissions (for HTTP)

#### Connection JSON file - SMPP address/port and HTTP URL

This is defined in JSON file. There may be multiple SMPP servers (addresses/ports) and multiple HTTP URLs.

#### Traffic template

This is Excel CSV file where each row describes exactly one SMS.

Excel is split in sections separated by empty row.

Each section has:

- Column caption row (one)
- SMS rows (one or many)
- Empty row

It is recommended that each section contains list of SMS for one
protocol. For example, first section SMPP, second HTTP, but format
described above allows template to be flexible and to use multiple
sections.

Each data row has always two columns:

1.  Timestamp (in RFC3339 format or some well-defined excel date format)
2.  Account name (or system ID if account is defined as SMPP)

##### SMPP columns

Besides first column that indicates timestamp and second that indicates
which account to use (system_id) other columns represent SMPP
submit_sm PDU fields.

- service_type string
- source_addr_ton int
- source_addr_npi int
- source_addr string
- dest_addr_ton int
- dest_addr_npi int
- destination_addr string
- esm_class int
- protocol_id int
- priority_flag int
- schedule_delivery_time string
- validity_period string
- registered_delivery int
- replace_if_present_flag int
- data_coding int
- sm_default_msg_id int
- sm_length int
- short_message string
- tlvs \[\]string

Each int parameter can be either decimal or hexadecimal number (hex is
indicated by prefix 0x, for example 0xe3. Each string parameter may be
given as string, with prefix s: (example "s:test" or as array of hex
bytes with prefix h (example h:31455e2043). When prefix is missing s: is
assumed.

TLV are special case, they are given as array of strings separated with
comma (",") and h: (hex format) is assumed.

##### HTTP columns

HTTP API is described on the following link:

<https://developers.horisen.com/en/sms-http-api>

Besides first column that indicates protocol and second that indicates
which account to use (auth.username) other represent fields in HTTP
protocol defined on above link:

- type string
- sender string
- receiver string
- dcs string
- text string
- dlrMask string

(DLR URL is generated by simulator itself).

##### Example of traffic template excel file

| Timestamp           | Account       | SMPP Fields… |
| ------------------- | ------------- | ------------ |
| 2021-05-11T11:12:34 | SMPP_ACCOUNT1 | …            |
| 2021-05-11T11:12:35 | SMPP_ACCOUNT2 | …            |
|                     |               |              |
| Timestamp           | Account       | HTTP Fields… |
| 2021-05-11T11:12:43 | HTTP_ACCOUNT1 | …            |
| 2021-05-11T11:12:43 | HTTP_ACCOUNT2 | …            |
| 2021-05-11T11:12:45 | HTTP_ACCOUNT3 | …            |
|                     |               |              |
| Timestamp           | Account       | SMPP Fields… |
| 2021-05-11T11:12:34 | SMPP_ACCOUNT3 | …            |
| 2021-05-11T11:12:35 | SMPP_ACCOUNT4 | …            |

##### Pattern variations

Pattern variation in a daily % range to change the values in logs and
stats:

###### Traffic Reduction Range (per day)

Example: 72-100% means, that randomly the simulator takes a value between 72 and 100 and reduced the traffic down to this level (reducing means: randomly skip messages to send in that amount) Example: Value picked today is 75%, this means that 25% of the SMS in pattern will be NOT sent, just ignored. Which one will not be sent, will be randomly

The value will be randomly picked on a daily basis for the whole traffic pattern to change the sending characteristics chosen.

###### Timestamp Shift Range (for every message)

Example: +/- 0.25-4h means, that randomly the simulator takes a value between +/- 0.25 and 4h and shift the timestamp in the traffic pattern accordingly. This will show in the statistics other sending patterns graphics.

The value will be randomly picked for every message in the traffic pattern to change the sending characteristics

###### Speed variation (per day)

Example: +/- 30-70% means, that randomly the simulator takes a value between +/- 30% and 70% and increases or decreases the speed (means adapt timestamp accordingly in traffic pattern) accordingly.

This will show in the statistics other peaks in statistic graphics. The value will be randomly picked on a daily basis for the whole traffic pattern to change the sending characteristics

Note: if the speed will be increased, and the pattern traffic has only 24h, it can be that at the end of the day some time has no traffic because the traffic pattern is used. In this case, just start the traffic pattern form the start again for the remaining time period.

###### Stress test mode

Ignore timestamps in traffic template (except for sorting messages), and send SMS with defined speed (SMS/sec).

#### How to simulate customer traffic

Each run of traffic template works in following manner. All the SMS are
sorted by timestamp and earliest timestamp is selected. All the other
timestamps are then calculated as "delay after" earliest.

All the SMPP account binds are established. If there is more than one
bind for some account, and there are multiple SMPP servers defined in
JSON, distribute connections evenly. If some of connection breaks, it
should be reconnected.

Simulator sends each SMS from template earliest after "delay" since
start of simulation. Following rules apply:

- If template defines that SMS is sent over SMPP account A, and
  account A has multiple binds, then SMS can be sent using any of
  binds belonging to account A.
- If template defines that SMS is sent over HTTP account B, any of
  URLs defined in connection JSON file can be used.

Traffic can be replayed with 2x, 3x or N times speed (meaning each delay
from start is reduced this many times) and speed is given as parameter
(command line or GUI).

### Supplier simulator details

Supplier simulator must implement server side of SMPP protocol and
server side of above HTTP protocol.

#### Defining supplier behavior

Its behavior is defined per username/system ID used for submission using
Excel CSV file with columns:

- Username/system_id
- Percentage of SMS to be accepted (in SMPP using submit_sm_resp, in HTTP using HTTP return code as described in protocol above). Rest are rejected.
- Percentage of SMS that are accepted that generates DELIVERED DLR
- Percentage of SMS that are accepted that generates UNDELIVERED DLR, with error code (may be multiple items for different error codes)
- Percentage of SMS that are accepted that generates REJECTED DLR (may be multiple items for different error codes)
- Percentage of SMS where DLR will not be generated

When generating DLR in the SMPP case, DLR can be sent over any bind
connected with the same account (system_id).

### Testing setup

There may be one or more testing setups. Testing setup is defined by
following data:

- set of ESME accounts §4.1,
- connection JSON file (that defines where to send traffic) §4.2
- supplier behavior §5.1

System should be able to work with multiple testing setups at the same time.

GUI/API must support defining Testing Setups – to import/export CSV/JSON
files and optionally to edit them directly in GUI.

### Testing session

Testing session is defined by following data

- selected testing setup
- traffic template §4.3
- pattern variations §4.3.4

Within one testing setup there may be multiple testing sessions, running concurrently.

Testing session can be run from GUI or by invoking API URL (e.g. from CRON).

GUI/API for managing sessions must include:

- Provisioning (defining new testing session, deleting, updating)
- Start testing session – start from the beginning
- Stop testing session – stop testing session
- Continue testing session – continue stopped testing session from the
  point where it stopped

### GUI requirements

GUI should provide authentication mechanism, and all APIs endpoints can
be invoked only when authenticated. Possible authentication mechanisms:

- OAuth2
- Basic HTTP authentication

### How to test

Connect customer simulator to supplier simulator and "replay" traffic
according to template.

### Some metrics about amount of data

- Up to 1000 SMPP and HTTP accounts
- Up to 10s of millions of messages in template.
