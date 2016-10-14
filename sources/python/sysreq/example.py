from canopsis.middleware.core import Middleware
from b3j0f.requester import Expression as E, Create as C, Read as R, Update as U, Delete as D

sysreq = Middleware.get_middleware_by_uri('sysreq://')

# first scenario: alarm stat

sysreq.driver.read(
    select=(E.Alarms.state, E.count(E.Alarms.id)),
    query=E.between(E.Alarms.ts, [E.now() - 24 * 3600, E.now()]) & (E.Alarms.data_id == E.Context.id) & E.in_('hg1', E.Context.hostgroups) & (E.Context.active == True),
    groupby='Alarms.state'
)


with sysreq.open() as transaction:

    transaction.read(
        select=E.count(E.Alarms),
        alias='alarmcount'
    )

    transaction.read(
        select=E.count(E.Alarms.id),
        query=E.Alarms.ack.ts > (E.Alarms.ts + (30 * 60)),
        alias='alarmnotsla'
    )

    transaction.read(
        select=E.alarmnotsla / E.alarmcount * 100,
        alias='alarmstat'
    )

# publish metric from alarm stat

with sysreq.open() as transaction:

    transaction.read(
        select=(E.count(E.Alarms.ack.ts > (E.Alarms.ts + (30 * 60))) / E.count(E.Alarms) * 100),
        alias='alarmstat'
    )

    transaction.create(
        E.Context,
        values={'type': 'metric', 'connector': 'canopsis', """..."""}
    )

    transaction.create(
        E.Perfdata,
        values={'value': E.alarmstat, 'name': 'alarmstat'}
    )


# with pbehavior (not downtime)
with sysreq.open() as transaction:

    transaction.read(
        select=(E.Alarms.state, E.count(E.Alarms.id)),
        query=(E.now() - 24 * 3600).as_('dtstart') & E.Pbehavior.whois(None, E.dtstart, E.now(), {'pbehavior': {'$neq': 'downtime'}}) & E.between(E.Alarms.ts, [E.now() - 24 * 3600, E.now()]) & (E.Alarms.data_id == E.Context.id) & E.in_('hg1', E.Context.hostgroups) & (E.Context.active == True),
        groupby='Alarms.state'
    )

    transaction.open().read(
        select=(E.Context.component.as_('metric'), E.avg(E.Perfdata.value).as_('value')),
        query=E.Perfdata.metric == 'sessionduration' & E.between(E.dtstart, E.now()) & E.Context.id == E.Perfdata.data_id,
        groupby='Context.component',
        alias='averages'
    ).commit()

    for perfdata in transaction.ctx['averages']:

        transaction.create(
            E.Perfdata,
            values=perfdata
        )
