# Deadline

# ---- Notice ----
This is a work in progress that is under heavy development and is under a lot of change. If it runs in an AT&T production environment, it does so under the supervision and maintance of it's maintainers. 

It is by no means ready for production use.

# Description

Deadline is an event driven schedule keeping application. Users define schedules to keep expecting events to be sent to the server.  It listens for these events, and if they do not happen by the appropritate time, the server will take some action as defined in the schedule (send an email for example). 

psuedo-code for a sample schedule may look like this:
* schedule. runs dayily. starts at 6:00 A.M.
	* 'wake-up' event expected by 7:00
	* 'showever' event expected by 7:30
	* 'leave-house' event expeted by 8:00
	* error handler is an email to jo424n@att.com 

### What deadline is
Deadline keeps a schedule 
* A simple way to keep a schedule based on the expectation that event *e* will occur by time *t*.

### What deadline is not
* A workflow management system. Deadline is a passive system, that is not actively triggering events or out of context 'things'. 


# Building
Run `make` or `make build` to make the binary. You'll need at least a go environment setup and gocylo in your $PATH
