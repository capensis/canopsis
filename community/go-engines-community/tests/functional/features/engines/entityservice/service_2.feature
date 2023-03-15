Feature: update service on event
  I need to be able to see new service state on event

  @concurrent
  Scenario: given entity service and ackremoved resource entity should update service counters on resource event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-second-1-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-second-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-service-second-1",
      "connector_name": "test-connector-name-service-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-second-1",
      "resource": "test-resource-service-second-1",
      "state": 1,
      "output": "test-output-service-second-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-second-1",
        "connector_name": "test-connector-name-service-second-1",
        "component": "test-component-service-second-1",
        "resource": "test-resource-service-second-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-service-second-1",
      "connector_name": "test-connector-name-service-second-1",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-service-second-1",
      "resource": "test-resource-service-second-1",
      "output": "test-output-service-second-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ack",
        "connector": "test-connector-service-second-1",
        "connector_name": "test-connector-name-service-second-1",
        "component": "test-component-service-second-1",
        "resource": "test-resource-service-second-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-service-second-1",
      "connector_name": "test-connector-name-service-second-1",
      "source_type": "resource",
      "event_type": "ackremove",
      "component": "test-component-service-second-1",
      "resource": "test-resource-service-second-1",
      "output": "test-output-service-second-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ackremove",
        "connector": "test-connector-service-second-1",
        "connector_name": "test-connector-name-service-second-1",
        "component": "test-component-service-second-1",
        "resource": "test-resource-service-second-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-service-second-1",
      "connector_name": "test-connector-name-service-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-second-1",
      "resource": "test-resource-service-second-1",
      "state": 3,
      "output": "test-output-service-second-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-service-second-1",
        "connector_name": "test-connector-name-service-second-1",
        "component": "test-component-service-second-1",
        "resource": "test-resource-service-second-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 3
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given entity service and updated resource entity by api should update service alarm on infos update
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-second-2-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "infos.manager",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-manager-service-second-2"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-second-2",
      "connector_name": "test-connector-name-service-second-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-second-2",
      "resource": "test-resource-service-second-2",
      "state": 3,
      "output": "test-output-service-second-2"
    }
    """
    When I do PUT /api/v4/entitybasics?_id=test-resource-service-second-2/test-component-service-second-2:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "infos": [
        {
          "description": "Manager",
          "name": "manager",
          "value": "test-manager-service-second-2"
        }
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entityupdated",
        "connector": "test-connector-service-second-2",
        "connector_name": "test-connector-name-service-second-2",
        "component": "test-component-service-second-2",
        "resource": "test-resource-service-second-2",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given deleted entity service should delete service alarm
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-second-3-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-second-3"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-service-second-3",
      "connector_name": "test-connector-name-service-second-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-second-3",
      "resource": "test-resource-service-second-3",
      "state": 3,
      "output": "test-output-service-second-3"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-second-3",
        "connector_name": "test-connector-name-service-second-3",
        "component": "test-component-service-second-3",
        "resource": "test-resource-service-second-3",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do DELETE /api/v4/entityservices/{{ .serviceID }}
    Then the response code should be 204
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "recomputeentityservice",
      "connector": "service",
      "connector_name": "service",
      "component": "{{ .serviceID }}",
      "source_type": "service"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-second-3",
      "connector_name": "test-connector-name-service-second-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-second-3",
      "resource": "test-resource-service-second-3",
      "state": 2,
      "output": "test-output-service-second-3"
    }
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }} until response code is 200 and body contains:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  @concurrent
  Scenario: given disabled entity service should not update service alarm
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-second-4-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-second-4"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-service-second-4",
      "connector_name": "test-connector-name-service-second-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-second-4",
      "resource": "test-resource-service-second-4",
      "state": 3,
      "output": "test-output-service-second-4"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-second-4",
        "connector_name": "test-connector-name-service-second-4",
        "component": "test-component-service-second-4",
        "resource": "test-resource-service-second-4",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-service-second-4-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": false,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-second-4"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "recomputeentityservice",
      "connector": "service",
      "connector_name": "service",
      "component": "{{ .serviceID }}",
      "source_type": "service"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-second-4",
      "connector_name": "test-connector-name-service-second-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-second-4",
      "resource": "test-resource-service-second-4",
      "state": 2,
      "output": "test-output-service-second-4"
    }
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }} until response code is 200 and body contains:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  @concurrent
  Scenario: given enabled entity service should create new service alarm
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-second-5-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-second-5"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-service-second-5",
      "connector_name": "test-connector-name-service-second-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-second-5",
      "resource": "test-resource-service-second-5",
      "state": 3,
      "output": "test-output-service-second-5"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-second-5",
        "connector_name": "test-connector-name-service-second-5",
        "component": "test-component-service-second-5",
        "resource": "test-resource-service-second-5",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-service-second-5-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": false,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-second-5"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    When I save response disableTimestamp={{ now }}
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "recomputeentityservice",
      "connector": "service",
      "connector_name": "service",
      "component": "{{ .serviceID }}",
      "source_type": "service"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-second-5",
      "connector_name": "test-connector-name-service-second-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-second-5",
      "resource": "test-resource-service-second-5",
      "state": 1,
      "output": "test-output-service-second-5"
    }
    """
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-service-second-5-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-second-5"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}&sort_by=v.resolved&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            }
          }
        },
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 1
            },
            "status": {
              "val": 1
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    When I save response serviceAlarm1ID={{ (index .lastResponse.data 0)._id }}
    When I save response serviceAlarm1Resolve={{ (index .lastResponse.data 0).v.resolved }}
    Then the difference between serviceAlarm1Resolve disableTimestamp is in range -2,2
    When I save response serviceAlarm2ID={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .serviceAlarm1ID }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      }
    ]
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .serviceAlarm2ID }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given deleted entity service should update impact service alarm
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-second-6-name-1",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-second-6"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response dependServiceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .dependServiceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .dependServiceID }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-second-6-name-2",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-service-second-6-name-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-service-second-6",
      "connector_name": "test-connector-name-service-second-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-second-6",
      "resource": "test-resource-service-second-6",
      "state": 3,
      "output": "test-output-service-second-6"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-second-6",
        "connector_name": "test-connector-name-service-second-6",
        "component": "test-component-service-second-6",
        "resource": "test-resource-service-second-6",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .dependServiceID }}"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do DELETE /api/v4/entityservices/{{ .dependServiceID }}
    Then the response code should be 204
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .dependServiceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 0
            },
            "status": {
              "val": 0
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "All: 0; Active: 0; Acknowledged: 0; NotAcknowledged: 0; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "service.service",
                "m": "All: 0; Active: 0; Acknowledged: 0; NotAcknowledged: 0; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 0
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given disabled entity service should update impact service alarm
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-second-7-name-1",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-second-7"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response dependServiceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .dependServiceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .dependServiceID }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-second-7-name-2",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-service-second-7-name-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-service-second-7",
      "connector_name": "test-connector-name-service-second-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-second-7",
      "resource": "test-resource-service-second-7",
      "state": 3,
      "output": "test-output-service-second-7"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-second-7",
        "connector_name": "test-connector-name-service-second-7",
        "component": "test-component-service-second-7",
        "resource": "test-resource-service-second-7",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .dependServiceID }}"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do PUT /api/v4/entityservices/{{ .dependServiceID }}:
    """json
    {
      "name": "test-entityservice-service-second-7-name-1",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": false,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-second-7"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .dependServiceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 0
            },
            "status": {
              "val": 0
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "All: 0; Active: 0; Acknowledged: 0; NotAcknowledged: 0; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "service.service",
                "m": "All: 0; Active: 0; Acknowledged: 0; NotAcknowledged: 0; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 0
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given disabled entity should update service alarm
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-second-8-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-second-8"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-service-second-8",
      "connector_name": "test-connector-name-service-second-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-second-8",
      "resource": "test-resource-service-second-8",
      "state": 3,
      "output": "test-output-service-second-8"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-second-8",
        "connector_name": "test-connector-name-service-second-8",
        "component": "test-component-service-second-8",
        "resource": "test-resource-service-second-8",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do PUT /api/v4/entitybasics?_id=test-resource-service-second-8/test-component-service-second-8:
    """json
    {
      "impact_level": 1,
      "enabled": false,
      "infos": [],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "connector": "test-connector-service-second-8",
        "connector_name": "test-connector-name-service-second-8",
        "component": "test-component-service-second-8",
        "resource": "test-resource-service-second-8",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 0
            },
            "status": {
              "val": 0
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "All: 0; Active: 0; Acknowledged: 0; NotAcknowledged: 0; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "service.service",
                "m": "All: 0; Active: 0; Acknowledged: 0; NotAcknowledged: 0; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 0
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given new entity service shouldn't count double ack
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-second-9",
      "connector_name": "test-connector-name-service-second-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-second-9",
      "resource": "test-resource-service-second-9",
      "state": 3,
      "output": "test-output-service-second-9"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-second-9-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-second-9"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "output": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-service-second-9",
      "connector_name": "test-connector-name-service-second-9",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-service-second-9",
      "resource": "test-resource-service-second-9",
      "output": "test-output-service-second-9"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ack",
        "connector": "test-connector-service-second-9",
        "connector_name": "test-connector-name-service-second-9",
        "component": "test-component-service-second-9",
        "resource": "test-resource-service-second-9",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-second-9",
      "connector_name": "test-connector-name-service-second-9",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-service-second-9",
      "resource": "test-resource-service-second-9",
      "output": "test-output-service-second-9"
    }
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "output": "All: 1; Active: 1; Acknowledged: 1; NotAcknowledged: 0; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
