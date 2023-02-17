Feature: update service on event
  I need to be able to see new service state on event

  @concurrent
  Scenario: given entity service and new resource entity should update service alarm on resource event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-1-name",
      "output_template": "Depends: {{ `{{ .Depends }}` }}; All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-1-1",
                "test-resource-service-1-2",
                "test-resource-service-1-3"
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
      "connector": "test-connector-service-1",
      "connector_name": "test-connector-name-service-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-1",
      "resource": "test-resource-service-1-1",
      "state": 1,
      "output": "test-output-service-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-1",
        "connector_name": "test-connector-name-service-1",
        "component": "test-component-service-1",
        "resource": "test-resource-service-1-1",
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
      "connector": "test-connector-service-1",
      "connector_name": "test-connector-name-service-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-1",
      "resource": "test-resource-service-1-2",
      "state": 3,
      "output": "test-output-service-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-1",
        "connector_name": "test-connector-name-service-1",
        "component": "test-component-service-1",
        "resource": "test-resource-service-1-2",
        "source_type": "resource"
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
      "connector": "test-connector-service-1",
      "connector_name": "test-connector-name-service-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-1",
      "resource": "test-resource-service-1-3",
      "state": 2,
      "output": "test-output-service-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-1",
        "connector_name": "test-connector-name-service-1",
        "component": "test-component-service-1",
        "resource": "test-resource-service-1-3",
        "source_type": "resource"
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
          "entity": {
            "_id": "{{ .serviceID }}",
            "depends_count": 3,
            "impacts_count": 0
          },
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "output": "Depends: 3; All: 3; Active: 3; Acknowledged: 0; NotAcknowledged: 3; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 1; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
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
                "m": "Depends: 1; All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "Depends: 1; All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "Depends: 2; All: 2; Active: 2; Acknowledged: 0; NotAcknowledged: 2; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
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
  Scenario: given entity service and new resource entity should update service alarm on service creation
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-2",
      "connector_name": "test-connector-name-service-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-2",
      "resource": "test-resource-service-2-1",
      "state": 3,
      "output": "test-output-service-2"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-2",
      "connector_name": "test-connector-name-service-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-2",
      "resource": "test-resource-service-2-2",
      "state": 2,
      "output": "test-output-service-2"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-2",
      "connector_name": "test-connector-name-service-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-2",
      "resource": "test-resource-service-2-3",
      "state": 1,
      "output": "test-output-service-2"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-2-name",
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
                "test-resource-service-2-1",
                "test-resource-service-2-2",
                "test-resource-service-2-3"
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
                "m": "All: 3; Active: 3; Acknowledged: 0; NotAcknowledged: 3; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 1; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 3; Active: 3; Acknowledged: 0; NotAcknowledged: 3; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 1; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
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
  Scenario: given entity service and removed resource entity should update service alarm on resource event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-3-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "infos.client",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-client-service-3"
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
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-3"
            }
          },
          {
            "field": "event_type",
            "cond": {
              "type": "eq",
              "value": "check"
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "client",
            "description": "Client",
            "value": "{{ `{{ .Event.ExtraInfos.client }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-eventfilter-service-3-description",
      "enabled": true,
      "priority": 2
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-service-3",
      "connector_name": "test-connector-name-service-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-3",
      "resource": "test-resource-service-3",
      "client": "test-client-service-3",
      "state": 2,
      "output": "test-output-service-3"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-3",
        "connector_name": "test-connector-name-service-3",
        "component": "test-component-service-3",
        "resource": "test-resource-service-3",
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
      "connector": "test-connector-service-3",
      "connector_name": "test-connector-name-service-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-3",
      "resource": "test-resource-service-3",
      "client": "test-another-client-service-3",
      "state": 2,
      "output": "test-output-service-3"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-service-3",
        "connector_name": "test-connector-name-service-3",
        "component": "test-component-service-3",
        "resource": "test-resource-service-3",
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
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
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
  Scenario: given entity service with updated pattern should update service alarm to increase state
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-4-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-4-1"
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
      "connector": "test-connector-service-4",
      "connector_name": "test-connector-name-service-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-4",
      "resource": "test-resource-service-4-1",
      "state": 2,
      "output": "test-output-service-4"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-4",
        "connector_name": "test-connector-name-service-4",
        "component": "test-component-service-4",
        "resource": "test-resource-service-4-1",
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-4",
      "connector_name": "test-connector-name-service-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-4",
      "resource": "test-resource-service-4-2",
      "state": 3,
      "output": "test-output-service-4"
    }
    """
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-service-4-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-4-2"
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
            "connector": "service",
            "connector_name": "service",
            "component": "{{ .serviceID }}",
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
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
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
  Scenario: given entity service with updated pattern should update service alarm to decrease state
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-5-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-5-1"
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
      "connector": "test-connector-service-5",
      "connector_name": "test-connector-name-service-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-5",
      "resource": "test-resource-service-5-1",
      "state": 2,
      "output": "test-output-service-5"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-5",
        "connector_name": "test-connector-name-service-5",
        "component": "test-component-service-5",
        "resource": "test-resource-service-5-1",
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-5",
      "connector_name": "test-connector-name-service-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-5",
      "resource": "test-resource-service-5-2",
      "state": 1,
      "output": "test-output-service-5"
    }
    """
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-service-5-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-5-2"
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
            "connector": "service",
            "connector_name": "service",
            "component": "{{ .serviceID }}",
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
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
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
  Scenario: given entity service and resolved resource entity should update service alarm on resource event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-6-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-6"
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
      "connector": "test-connector-service-6",
      "connector_name": "test-connector-name-service-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-6",
      "resource": "test-resource-service-6",
      "state": 2,
      "output": "test-output-service-6"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-6",
        "connector_name": "test-connector-name-service-6",
        "component": "test-component-service-6",
        "resource": "test-resource-service-6",
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-6",
      "connector_name": "test-connector-name-service-6",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-service-6",
      "resource": "test-resource-service-6",
      "output": "test-output-service-6"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-service-6",
      "connector_name": "test-connector-name-service-6",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-service-6",
      "resource": "test-resource-service-6",
      "output": "test-output-service-6"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "resolve_cancel",
        "connector": "test-connector-service-6",
        "connector_name": "test-connector-name-service-6",
        "component": "test-component-service-6",
        "resource": "test-resource-service-6",
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
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
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
  Scenario: given new entity service and resolved resource entity should not create service alarm on service creation
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-7",
      "connector_name": "test-connector-name-service-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-7",
      "resource": "test-resource-service-7",
      "state": 2,
      "output": "test-output-service-7"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-7",
      "connector_name": "test-connector-name-service-7",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-service-7",
      "resource": "test-resource-service-7",
      "output": "test-output-service-7"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-7",
      "connector_name": "test-connector-name-service-7",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-service-7",
      "resource": "test-resource-service-7",
      "output": "test-output-service-7"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-7-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-7"
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
        "component": "{{ .serviceID }}"
      }
    ]
    """
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

  @concurrent
  Scenario: given new entity service should count all alarms on service creation
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-8",
      "connector_name": "test-connector-name-service-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-8",
      "resource": "test-resource-service-8-1",
      "state": 1,
      "output": "test-output-service-8"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-8",
      "connector_name": "test-connector-name-service-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-8",
      "resource": "test-resource-service-8-2",
      "state": 2,
      "output": "test-output-service-8"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-8",
      "connector_name": "test-connector-name-service-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-8",
      "resource": "test-resource-service-8-3",
      "state": 3,
      "output": "test-output-service-8"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-8",
      "connector_name": "test-connector-name-service-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-8",
      "resource": "test-resource-service-8-4",
      "state": 0,
      "output": "test-output-service-8"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-8",
      "connector_name": "test-connector-name-service-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-8",
      "resource": "test-resource-service-8-5",
      "state": 3,
      "output": "test-output-service-8"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-8",
      "connector_name": "test-connector-name-service-8",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-service-8",
      "resource": "test-resource-service-8-5",
      "output": "test-output-service-8"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-8",
      "connector_name": "test-connector-name-service-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-8",
      "resource": "test-resource-service-8-6",
      "state": 3,
      "output": "test-output-service-8"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-service-8",
      "connector_name": "test-connector-name-service-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-8",
      "resource": "test-resource-service-8-6",
      "state": 0,
      "output": "test-output-service-8"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-8-name",
      "output_template": "Depends: {{ `{{ .Depends}}` }}; All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-8-1",
                "test-resource-service-8-2",
                "test-resource-service-8-3",
                "test-resource-service-8-4",
                "test-resource-service-8-5",
                "test-resource-service-8-6"
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
            "output": "Depends: 6; All: 5; Active: 5; Acknowledged: 1; NotAcknowledged: 4; AcknowledgedUnderPbh: 0; StateCritical: 2; StateMajor: 1; StateMinor: 1; StateOk: 1; Pbehaviors: map[]; UnderPbehavior: 0;",
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
                "m": "Depends: 6; All: 5; Active: 5; Acknowledged: 1; NotAcknowledged: 4; AcknowledgedUnderPbh: 0; StateCritical: 2; StateMajor: 1; StateMinor: 1; StateOk: 1; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "Depends: 6; All: 5; Active: 5; Acknowledged: 1; NotAcknowledged: 4; AcknowledgedUnderPbh: 0; StateCritical: 2; StateMajor: 1; StateMinor: 1; StateOk: 1; Pbehaviors: map[]; UnderPbehavior: 0;",
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
  Scenario: given entity service and acked resource entity should update service counters on resource event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-9-name",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-9"
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
      "connector": "test-connector-service-9",
      "connector_name": "test-connector-name-service-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-9",
      "resource": "test-resource-service-9",
      "state": 1,
      "output": "test-output-service-9"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-9",
        "connector_name": "test-connector-name-service-9",
        "component": "test-component-service-9",
        "resource": "test-resource-service-9",
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
      "connector": "test-connector-service-9",
      "connector_name": "test-connector-name-service-9",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-service-9",
      "resource": "test-resource-service-9",
      "output": "test-output-service-9"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ack",
        "connector": "test-connector-service-9",
        "connector_name": "test-connector-name-service-9",
        "component": "test-component-service-9",
        "resource": "test-resource-service-9",
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
      "connector": "test-connector-service-9",
      "connector_name": "test-connector-name-service-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-9",
      "resource": "test-resource-service-9",
      "state": 3,
      "output": "test-output-service-9"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-service-9",
        "connector_name": "test-connector-name-service-9",
        "component": "test-component-service-9",
        "resource": "test-resource-service-9",
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
                "m": "All: 1; Active: 1; Acknowledged: 1; NotAcknowledged: 0; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
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
