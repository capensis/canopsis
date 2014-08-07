component = "test_event_32"
        # status legend:
        # 0 == Ok
        # 1 == On going
        # 2 == Stealthy
        # 3 == Bagot
        # 4 == Canceled


def e1(previousEvent):

    return ('event l1', {
        "domain": "Infra",
        "perimeter": "Iris",
        "connector": "nagios",
        "connector_name": "nagios1",
        "event_type": "check",
        "source_type": "resource",
        "component": "serv1",
        "resource": "Disk",
        "state": 0,
        "status": 1,
        'criticity': 1,
        "state_type": 1,
        'hostgroups': ['HG3', 'HG4'],
        "output": "",
        "display_name": "DISPLAY_NAME"
    })

def e2(previousEvent):

    return ('event l2', {
        "domain": "Infra",
        "perimeter": "Iris",
        "connector": "shinken",
        "connector_name": "nagios1",
        "event_type": "check",
        "source_type": "resource",
        "component": "serv1",
        "resource": "RAM",
        'criticity': 2,
        "state": 1,
        "status": 3,
        "state_type": 1,
        'hostgroups': ['HG3', 'HG4']
    })

def e3(previousEvent):

    return ('event l3', {
        "domain": "Appli",
        "perimeter": "Ravel",
        "connector": "nagios",
        "connector_name": "nagios1",
        "event_type": "check",
        "source_type": "resource",
        "component": "test",
        "resource": "VM",
        'criticity': 3,
        "state": 1,
        "status": 2,
        "state_type": 1,
        "output": "/usr/local/nagios/libexec/check_ping -H 192.168.1.2 -w 100.0,90% -c 200.0,60%",
        'hostgroups': ['HG3', 'HG4']
    })

def e3(previousEvent):

    return ('event OK', {
        "domain": "Appli",
        "perimeter": "Trans",
        "connector": "nagios",
        "connector_name": "nagios1",
        "event_type": "check",
        "source_type": "resource",
        "component": "Prod",
        "resource": "VM",
        'criticity': 4,
        "state": 1,
        "status": 2,
        "state_type": 0,
        "output": "/usr/local/nagios/libexec/check_ping -H 192.168.1.2 -w 100.0,90% -c 200.0,60%",
        'hostgroups': ['HG3', 'HG4']
    })

def e4(previousEvent):

    return ('event l4', {
        "domain": "Appli",
        "perimeter": "Trans",
        "connector": "nagios",
        "connector_name": "nagios1",
        "event_type": "check",
        "source_type": "resource",
        "component": "Prod",
        "resource": "VM",
        "state": 3,
        'criticity': 1,
        "status": 2,
        "state_type": 0,
        "output": "/usr/local/nagios/libexec/check_ping -H 192.168.1.2 -w 100.0,90% -c 200.0,60%",
        'hostgroups': ['HG3', 'HG4']
    })

def e5(previousEvent):

    return ('event l5', {
        "domain": "Infra",
        "perimeter": "RAVEL",
        "connector": "schneider",
        "connector_name": "nagios1",
        "event_type": "check",
        "source_type": "resource",
        "component": "serv2",
        "resource": "Disk",
        "state": 2,
        "status": 2,
        "state_type": 0,
        "output": "",
        'hostgroups': ['HG3', 'HG4']
    })


def e6(previousEvent):

    return ('event l6', {
        "domain": "Infra",
        "perimeter": "RAVEL",
        "connector": "shinken",
        "connector_name": "nagios1",
        "event_type": "check",
        "source_type": "resource",
        "component": "serv2",
        "resource": "Ram",
        "state": 3,
        "status": 0,
        "state_type": 0,
        "output": "",
        'hostgroups': ['HG3', 'HG4']
    })

scenario = [e1, e2, e3, e4, e5]
