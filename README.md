# Deadline

# ---- **Notice** ----
This is a work in progress. If it runs in an AT&T environment, it does so under the supervision of it's maintainers.

It is by no means ready for production use.

# Description

Deadline is an application for keeping track of schedules. In distributed systems it's often very hard to keep track of data.  Especially as it changes over time and by different systems (often different teams).

Deadline is not a workflow management system.  It is more accurately described as a workflow *observation* system.  It is strictly passive with regard to the actual events and workflows it's observing.

We recognize that workflow management is increasingly difficult in polyglot enterprises and is often hard (sometimes impossible) to unify under one master workflow management system.  Thus we offer third party observability (and alerting), because it is significantly cheaper in terms of integration work yet still offers a lot value. Especially with regard to troubleshooting and reporting.

# Schedules

Schedules are structured as [Directed Acyclic Graphs](https://en.wikipedia.org/wiki/Directed_acyclic_graph). Most nodes are `event` nodes that have constraints on them.  Typically it is the expectations that event *e* occurs by time *t*, though you could have further constraints, for example, on the content body of the event.  `Handler` nodes then handle the failure when this expectation is not met, usually in some form of an alert.  Other node types are control nodes like `start`, `end`.

This example is in fact a re-worded test cases, which you can find [here](/dao/testdata/single_event_schedule.xml).  It's the simplest of cases with just one event.  The use case it depicts is this: I want to be sure that system X generates a report every day.  If it doesn't happen in a timely manner (say within 3 hours of the expected generation time), I want an email sent to me.

#### Schedule start
```xml
<schedule timing="24h" name="my_daily_report" starts-at="2018-09-08T00:00:00Z">
    <start to="sys_x_generated_report" />
```
Looking at the `schedule` element we can see that the `timing` is 24h, meaning it recurs every 24 hours, and `starts-at` midnight.  The `name` of the schedule is my_daily_report.

Every schedule must have at least a `start` and `end` node. The start node is the entry point of the schedule and it's only functionality is to designate where in the DAG to go at the very beginning of the schedule

In this case we want to go `to` (the xml attribute) a node named sys_x_generated_report.

#### Events
```xml
    <event name="sys_x_generated_report" ok="scheduleEnd" error="email error">
        <constraints>
            <receive-by>3h</receive-by>
        </constraints>
    </event>
```
Every event node must have `name`, `ok` and `error` attributes.  If an event meets it's constraints, the next step in the schedule is the node defined in the ok attribute.  If, however, the event fails to meet it's constraints the schedule's next step is to the node defined in the error attribute.

Event nodes hold `constraints` which is where a lot deadline's functionality comes from.  This is a constraint to receive this event (`receive-by` element) 3 hours after the start of the schedule.  Recieive by constraints are always relative to the start time of the schedule.

#### Handlers
```xml
<handler name="email error" to="scheduleEnd">
	<email>
		<to>me@mycompany.com</to>
		<message>That report we need didn't generate today.</message>
		<subject>report generation failed</subject>
	</email>
</handler>
```

`Handler` nodes are the other half deadlines core functionality.  They handle failures by doing things like sending emails.  In this case, and `email` element is defined to email `to` me@mycompany.com that the report was not generated.

Note that this node goes `to` the end node. Handlers could alternatively go to any other sort of node. 

#### Schedule ends
```xml
	<end name="scheduleEnd" />
</schedule>
```

End elements simply need a `name` and they denote the end of the schedule.

#### The full example

```xml
<schedule timing="daily" name="my_daily_report" starts-at="2018-09-08T00:00:00Z">
    <start to="sys_x_generated_report" />

    <event name="sys_x_generated_report" ok="scheduleEnd" error="email error">
        <constraints>
            <receive-by>3h</receive-by>
        </constraints>
    </event>

    <handler name="email error" to="scheduleEnd">
        <email>
            <to>me@mycompany.com</to>
	    <message>That report we need didn't generate today.</message>
	    <subject>report generation failed</subject>
        </email>
    </handler>

    <end name="scheduleEnd" />
</schedule>
```

# Collecting events

Deadline currently only accepts 'events' through an HTTP interface detailed below.

## HTTP Messages

Events can be collected through the `/api/v1/event` API. The method is PUT and you'll also need to specify the content-type as application/json.

Here's an example you may send.  Only `name` is required and it's a string.  `details` is a json object where you can specify any thing you like.
```json
{
	"name": "that event you wanted",
	"details": {
		"from": "system x",
		"the_answer": 42,
		"success": true,
		"other_object": {}
	},
}
```


# Design goals

* Deadline is a strictly passive listener of events.  Although handlers could be used to modify external systems (i.e., a webhook handler that triggers a retry on an external system), this pattern should be avoided.  Deadline strives to be a workflow *observation* system which is to say users should treat the workflows it's observing as read only.
* Schedules have an explicit start times re-occurrence intervals (one-time occurrence being a special case). Deadline will only evaluate one occurrence of a schedule at a time.
* Deadline is meant for schedules that last hours and events that are occur minutes or hours apart rather than seconds or milliseconds.
* Deadline should collect events from as many interfaces as applicable and beneficial. In the least it exposes an HTTP interface but should be extended to be a subscriber to messaging systems.
* Handlers, by definition '*do something*' so they should support many things to do like web-hooks
* Ecosystem of *light* integration points.  *Light* is meant to denote an integration point that does not impose a dependency on Deadline. For example, log parsers that search for a particular message or list of regular expressions and then sends an event to Deadline.  This is an out of procress integration that does not affect the system writing the log.

# Building
Run `make` or `make build` to make the binary. You'll need at least a go environment setup and gocylo in your $PATH
