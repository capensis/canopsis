Feature: no update service when entity is inactive
  I need to be able to not update service when pause or maintenance pbehavior is in action.

  Scenario: given entity service and maintenance pbehavior should not update service alarm on create with pbhenter event
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-service-1",
      "connector_name": "test-connector-name-pbehavior-service-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-service-1",
      "resource": "test-resource-pbehavior-service-1",
      "state": 0,
      "output": "test-output-pbehavior-service-1"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-pbehavior-service-1-name",
      "output_template": "Depends: {{ `{{ .Depends }}` }}; All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acked: {{ `{{.Acknowledged}}` }}; NotAcked: {{ `{{.NotAcknowledged}}` }}; AckedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-service-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-service-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-service-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-service-1",
      "connector_name": "test-connector-name-pbehavior-service-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-service-1",
      "resource": "test-resource-pbehavior-service-1",
      "state": 1,
      "output": "test-output-pbehavior-service-1"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
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

  Scenario: given entity service and maintenance pbehavior should update service alarm on pbhenter event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-pbehavior-service-2-name",
      "output_template": "Depends: {{ `{{ .Depends }}` }}; All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acked: {{ `{{.Acknowledged}}` }}; NotAcked: {{ `{{.NotAcknowledged}}` }}; AckedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-service-2"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-service-2",
      "connector_name": "test-connector-name-pbehavior-service-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-service-2",
      "resource": "test-resource-pbehavior-service-2",
      "state": 3,
      "output": "test-output-pbehavior-service-2"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-service-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-service-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
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
                "m": "Depends: 1; All: 1; Active: 1; Acked: 0; NotAcked: 1; AckedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Active: 1; Acked: 0; NotAcked: 1; AckedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Active: 0; Acked: 0; NotAcked: 0; AckedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateOk: 1; Pbehaviors: map[test-maintenance-type-to-engine:1]; UnderPbehavior: 1;",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Active: 0; Acked: 0; NotAcked: 0; AckedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateOk: 1; Pbehaviors: map[test-maintenance-type-to-engine:1]; UnderPbehavior: 1;",
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

  Scenario: given entity service and maintenance pbehavior should update service alarm on pbhleave event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-pbehavior-service-3-name",
      "output_template": "Depends: {{ `{{ .Depends }}` }}; All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acked: {{ `{{.Acknowledged}}` }}; NotAcked: {{ `{{.NotAcknowledged}}` }}; AckedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-service-3"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-service-3",
      "connector_name": "test-connector-name-pbehavior-service-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-service-3",
      "resource": "test-resource-pbehavior-service-3",
      "state": 3,
      "output": "test-output-pbehavior-service-3"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-service-3",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2s" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-service-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 4 events processing
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
                "m": "Depends: 1; All: 1; Active: 1; Acked: 0; NotAcked: 1; AckedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Active: 1; Acked: 0; NotAcked: 1; AckedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Active: 0; Acked: 0; NotAcked: 0; AckedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateOk: 1; Pbehaviors: map[test-maintenance-type-to-engine:1]; UnderPbehavior: 1;",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Active: 0; Acked: 0; NotAcked: 0; AckedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateOk: 1; Pbehaviors: map[test-maintenance-type-to-engine:1]; UnderPbehavior: 1;",
                "val": 0
              },
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Active: 1; Acked: 0; NotAcked: 1; AckedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[test-maintenance-type-to-engine:0]; UnderPbehavior: 0;",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Active: 1; Acked: 0; NotAcked: 1; AckedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[test-maintenance-type-to-engine:0]; UnderPbehavior: 0;",
                "val": 1
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 6
            }
          }
        }
      }
    ]
    """

  Scenario: given entity service and maintenance and inactive pbehaviors should update service alarm on pbhleave and enter event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-pbehavior-service-4-name",
      "output_template": "Depends: {{ `{{ .Depends }}` }}; All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acked: {{ `{{.Acknowledged}}` }}; NotAcked: {{ `{{.NotAcknowledged}}` }}; AckedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-pbehavior-service-4-1",
                "test-resource-pbehavior-service-4-2"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-service-4",
      "connector_name": "test-connector-name-pbehavior-service-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-service-4",
      "resource": "test-resource-pbehavior-service-4-1",
      "state": 3,
      "output": "test-output-pbehavior-service-4"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-service-4-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2s" }},
      "color": "#FFFFFF",
      "type": "test-inactive-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-service-4-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-service-4-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-service-4-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 4 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-service-4",
      "connector_name": "test-connector-name-pbehavior-service-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-service-4",
      "resource": "test-resource-pbehavior-service-4-2",
      "state": 2,
      "output": "test-output-pbehavior-service-4"
    }
    """
    When I wait the end of 2 events processing
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
              "val": 2
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
                "m": "Depends: 1; All: 1; Active: 1; Acked: 0; NotAcked: 1; AckedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Active: 1; Acked: 0; NotAcked: 1; AckedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Active: 0; Acked: 0; NotAcked: 0; AckedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateOk: 1; Pbehaviors: map[test-inactive-type-to-engine:1]; UnderPbehavior: 1;",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Active: 0; Acked: 0; NotAcked: 0; AckedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateOk: 1; Pbehaviors: map[test-inactive-type-to-engine:1]; UnderPbehavior: 1;",
                "val": 0
              },
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "Depends: 2; All: 2; Active: 1; Acked: 0; NotAcked: 1; AckedUnderPbh: 0; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateOk: 1; Pbehaviors: map[test-inactive-type-to-engine:0 test-maintenance-type-to-engine:1]; UnderPbehavior: 1;",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "Depends: 2; All: 2; Active: 1; Acked: 0; NotAcked: 1; AckedUnderPbh: 0; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateOk: 1; Pbehaviors: map[test-inactive-type-to-engine:0 test-maintenance-type-to-engine:1]; UnderPbehavior: 1;",
                "val": 1
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 6
            }
          }
        }
      }
    ]
    """
