Feature: update service on event
  I need to be able to see new service state on event

  @concurrent
  Scenario: given new entity service with dependencies and bulk enable and bulk disable requests should recompute service
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-third-1",
      "connector_name": "test-connector-name-service-third-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-third-1",
      "resource": "test-resource-service-third-1-1",
      "state": 1,
      "output": "test-output-service-third-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-third-1",
      "connector_name": "test-connector-name-service-third-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-third-1",
      "resource": "test-resource-service-third-1-2",
      "state": 2,
      "output": "test-output-service-third-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-third-1",
      "connector_name": "test-connector-name-service-third-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-third-1",
      "resource": "test-resource-service-third-1-3",
      "state": 3,
      "output": "test-output-service-third-1"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-third-1-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-third-1-1",
                "test-resource-service-third-1-2",
                "test-resource-service-third-1-3"
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
            "output": "All: 3; Active: 3; Acknowledged: 0; NotAcknowledged: 3; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 1; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;"
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
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "test-resource-service-third-1-2/test-component-service-third-1"
      },
      {
        "_id": "test-resource-service-third-1-3/test-component-service-third-1"
      }
    ]
    """
    Then the response code should be 207
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "connector": "test-connector-service-third-1",
        "connector_name": "test-connector-name-service-third-1",
        "component": "test-component-service-third-1",
        "resource": "test-resource-service-third-1-2",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "connector": "test-connector-service-third-1",
        "connector_name": "test-connector-name-service-third-1",
        "component": "test-component-service-third-1",
        "resource": "test-resource-service-third-1-3",
        "source_type": "resource"
      },
      {
        "event_type": "check",
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
              "val": 1
            },
            "status": {
              "val": 1
            },
            "output": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;"
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
    When I do PUT /api/v4/bulk/entities/enable:
    """json
    [
      {
        "_id": "test-resource-service-third-1-2/test-component-service-third-1"
      },
      {
        "_id": "test-resource-service-third-1-3/test-component-service-third-1"
      }
    ]
    """
    Then the response code should be 207
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "connector": "test-connector-service-third-1",
        "connector_name": "test-connector-name-service-third-1",
        "component": "test-component-service-third-1",
        "resource": "test-resource-service-third-1-2",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "connector": "test-connector-service-third-1",
        "connector_name": "test-connector-name-service-third-1",
        "component": "test-component-service-third-1",
        "resource": "test-resource-service-third-1-3",
        "source_type": "resource"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-service-third-1",
      "connector_name": "test-connector-name-service-third-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-third-1",
      "resource": "test-resource-service-third-1-2",
      "state": 2,
      "output": "test-output-service-third-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-third-1",
        "connector_name": "test-connector-name-service-third-1",
        "component": "test-component-service-third-1",
        "resource": "test-resource-service-third-1-2",
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
      "connector": "test-connector-service-third-1",
      "connector_name": "test-connector-name-service-third-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-third-1",
      "resource": "test-resource-service-third-1-3",
      "state": 3,
      "output": "test-output-service-third-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-third-1",
        "connector_name": "test-connector-name-service-third-1",
        "component": "test-component-service-third-1",
        "resource": "test-resource-service-third-1-3",
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
            },
            "output": "All: 3; Active: 3; Acknowledged: 0; NotAcknowledged: 3; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 1; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;"
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

  @concurrent
  Scenario: given new entity service with services as dependencies and bulk enable and bulk disable requests should recompute service
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-third-2",
      "connector_name": "test-connector-name-service-third-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-third-2",
      "resource": "test-resource-service-third-2-1",
      "state": 1,
      "output": "test-output-service-third-2"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-third-2",
      "connector_name": "test-connector-name-service-third-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-third-2",
      "resource": "test-resource-service-third-2-2",
      "state": 2,
      "output": "test-output-service-third-2"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-third-2",
      "connector_name": "test-connector-name-service-third-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-third-2",
      "resource": "test-resource-service-third-2-3",
      "state": 3,
      "output": "test-output-service-third-2"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-third-2-name-1",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-third-2-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID1={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID1 }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID1 }}"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-third-2-name-2",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-third-2-2"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID2={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID2 }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID2 }}"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-third-2-name-3",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-third-2-3"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID3={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID3 }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID3 }}"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-third-2-name-4",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-entityservice-service-third-2-name-1",
                "test-entityservice-service-third-2-name-2",
                "test-entityservice-service-third-2-name-3"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response impactServiceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .impactServiceID }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .impactServiceID }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .impactServiceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .impactServiceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "output": "All: 3; Active: 3; Acknowledged: 0; NotAcknowledged: 3; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 1; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;"
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
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "{{ .serviceID2 }}"
      },
      {
        "_id": "{{ .serviceID3 }}"
      }
    ]
    """
    Then the response code should be 207
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID2 }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .impactServiceID }}",
        "source_type": "service"
      },
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID3 }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .impactServiceID }}",
        "source_type": "service"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .impactServiceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .impactServiceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 1
            },
            "status": {
              "val": 1
            },
            "output": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;"
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
    When I do PUT /api/v4/bulk/entities/enable:
    """json
    [
      {
        "_id": "{{ .serviceID2 }}"
      },
      {
        "_id": "{{ .serviceID3 }}"
      }
    ]
    """
    Then the response code should be 207
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID2 }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID2 }}"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .impactServiceID }}",
        "source_type": "service"
      },
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID3 }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID3 }}"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .impactServiceID }}",
        "source_type": "service"
      }
    ]
    """
    When I do GET /api/v4/alarms?search={{ .impactServiceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .impactServiceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "output": "All: 3; Active: 3; Acknowledged: 0; NotAcknowledged: 3; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 1; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;"
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
