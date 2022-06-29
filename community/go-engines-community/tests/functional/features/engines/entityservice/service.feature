Feature: update service on event
  I need to be able to see new service state on event

  Scenario: given entity service and new resource entity should update service alarm on resource event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-1-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
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
    When I wait the end of 2 events processing
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
    When I wait the end of 2 events processing
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
    When I wait the end of 2 events processing
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 2; Alarms: 2; Acknowledged: 0; NotAcknowledged: 2; StateCritical: 1; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given entity service and new resource entity should update service alarm on service creation
    Given I am admin
    When I send an event:
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
    When I wait the end of event processing
    When I send an event:
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
    When I wait the end of event processing
    When I send an event:
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
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-2-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
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
                "m": "All: 3; Alarms: 3; Acknowledged: 0; NotAcknowledged: 3; StateCritical: 1; StateMajor: 1; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 3; Alarms: 3; Acknowledged: 0; NotAcknowledged: 3; StateCritical: 1; StateMajor: 1; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given entity service and removed resource entity should update service alarm on resource event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-3-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
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
    When I wait the end of 2 events processing
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
      "actions": [
        {
          "type": "set_entity_info_from_template",
          "name": "client",
          "description": "Client",
          "value": "{{ `{{ .Event.ExtraInfos.client }}` }}"
        }
      ],
      "description": "test-eventfilter-service-3-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
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
    When I wait the end of 2 events processing
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "All: 0; Alarms: 0; Acknowledged: 0; NotAcknowledged: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "service.service",
                "m": "All: 0; Alarms: 0; Acknowledged: 0; NotAcknowledged: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given entity service with updated pattern should update service alarm to increase state
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-4-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
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
    When I wait the end of 2 events processing
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
    When I wait the end of 2 events processing
    When I send an event:
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
    When I wait the end of event processing
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-service-4-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
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
    When I wait the end of 2 events processing
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given entity service with updated pattern should update service alarm to decrease state
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-5-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
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
    When I wait the end of 2 events processing
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
    When I wait the end of 2 events processing
    When I send an event:
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
    When I wait the end of event processing
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-service-5-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
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
    When I wait the end of 2 events processing
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given entity service and resolved resource entity should update service alarm on resource event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-6-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
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
    When I wait the end of 2 events processing
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
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-6",
      "connector_name": "test-connector-name-service-6",
      "source_type": "resource",
      "event_type": "done",
      "component": "test-component-service-6",
      "resource": "test-resource-service-6",
      "output": "test-output-service-6"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-6",
      "connector_name": "test-connector-name-service-6",
      "source_type": "resource",
      "event_type": "resolve_done",
      "component": "test-component-service-6",
      "resource": "test-resource-service-6",
      "output": "test-output-service-6"
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 1; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "All: 0; Alarms: 0; Acknowledged: 0; NotAcknowledged: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "service.service",
                "m": "All: 0; Alarms: 0; Acknowledged: 0; NotAcknowledged: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given new entity service and resolved resource entity should not create service alarm on service creation
    Given I am admin
    When I send an event:
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
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-7",
      "connector_name": "test-connector-name-service-7",
      "source_type": "resource",
      "event_type": "done",
      "component": "test-component-service-7",
      "resource": "test-resource-service-7",
      "output": "test-output-service-7"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-7",
      "connector_name": "test-connector-name-service-7",
      "source_type": "resource",
      "event_type": "resolve_done",
      "component": "test-component-service-7",
      "resource": "test-resource-service-7",
      "output": "test-output-service-7"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-7-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
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

  Scenario: given new entity service should count all alarms on service creation
    Given I am admin
    When I send an event:
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
    When I wait the end of event processing
    When I send an event:
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
    When I wait the end of event processing
    When I send an event:
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
    When I wait the end of event processing
    When I send an event:
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
    When I wait the end of event processing
    When I send an event:
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
    When I wait the end of event processing
    When I send an event:
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
    When I wait the end of event processing
    When I send an event:
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
    When I wait the end of event processing
    When I send an event:
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
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-8-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
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
                "m": "All: 5; Alarms: 5; Acknowledged: 1; NotAcknowledged: 4; StateCritical: 2; StateMajor: 1; StateMinor: 1; StateInfo: 1; Pbehaviors: map[];",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 5; Alarms: 5; Acknowledged: 1; NotAcknowledged: 4; StateCritical: 2; StateMajor: 1; StateMinor: 1; StateInfo: 1; Pbehaviors: map[];",
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

  Scenario: given entity service and acked resource entity should update service counters on resource event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-9-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
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
    When I wait the end of 2 events processing
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
    When I wait the end of 2 events processing
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
    When I wait the end of 2 events processing
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 1; NotAcknowledged: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given entity service and ackremoved resource entity should update service counters on resource event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-10-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-10"
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
      "connector": "test-connector-service-10",
      "connector_name": "test-connector-name-service-10",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-10",
      "resource": "test-resource-service-10",
      "state": 1,
      "output": "test-output-service-10"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-10",
      "connector_name": "test-connector-name-service-10",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-service-10",
      "resource": "test-resource-service-10",
      "output": "test-output-service-10"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-10",
      "connector_name": "test-connector-name-service-10",
      "source_type": "resource",
      "event_type": "ackremove",
      "component": "test-component-service-10",
      "resource": "test-resource-service-10",
      "output": "test-output-service-10"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-10",
      "connector_name": "test-connector-name-service-10",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-10",
      "resource": "test-resource-service-10",
      "state": 3,
      "output": "test-output-service-10"
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given entity service and new resource entity with component infos by extra infos should update service alarm on resource event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-11-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "component_infos.client",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-client-service-11"
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
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-service-11"
            }
          },
          {
            "field": "source_type",
            "cond": {
              "type": "eq",
              "value": "component"
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
      "actions": [
        {
          "type": "set_entity_info_from_template",
          "name": "client",
          "description": "Client",
          "value": "{{ `{{ .Event.ExtraInfos.client }}` }}"
        }
      ],
      "description": "test-eventfilter-service-11-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-service-11",
      "connector_name": "test-connector-name-service-11",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-service-11",
      "state": 1,
      "output": "test-output-service-11",
      "client": "test-client-service-11"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-11",
      "connector_name": "test-connector-name-service-11",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-11",
      "resource": "test-resource-service-11",
      "state": 3,
      "output": "test-output-service-11"
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given entity service and new resource entity with component infos by extra infos should update service alarm on component event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-12-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "component_infos.client",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-client-service-12"
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
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-service-12"
            }
          },
          {
            "field": "source_type",
            "cond": {
              "type": "eq",
              "value": "component"
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
      "actions": [
        {
          "type": "set_entity_info_from_template",
          "name": "client",
          "description": "Client",
          "value": "{{ `{{ .Event.ExtraInfos.client }}` }}"
        }
      ],
      "description": "test-eventfilter-service-12-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-service-12",
      "connector_name": "test-connector-name-service-12",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-12",
      "resource": "test-resource-service-12",
      "state": 3,
      "output": "test-output-service-12"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-12",
      "connector_name": "test-connector-name-service-12",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-service-12",
      "state": 1,
      "output": "test-output-service-12",
      "client": "test-client-service-12"
    }
    """
    When I wait the end of 3 events processing
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given entity service and removed resource entity with component infos by extra infos should update service alarm on component event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-13-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "component_infos.client",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-client-service-13"
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
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-service-13"
            }
          },
          {
            "field": "source_type",
            "cond": {
              "type": "eq",
              "value": "component"
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
      "actions": [
        {
          "type": "set_entity_info_from_template",
          "name": "client",
          "description": "Client",
          "value": "{{ `{{ .Event.ExtraInfos.client }}` }}"
        }
      ],
      "description": "test-eventfilter-service-13-description",
      "enabled": true,
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass"
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-service-13",
      "connector_name": "test-connector-name-service-13",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-13",
      "resource": "test-resource-service-13",
      "state": 3,
      "output": "test-output-service-13"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-13",
      "connector_name": "test-connector-name-service-13",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-service-13",
      "state": 1,
      "output": "test-output-service-13",
      "client": "test-client-service-13"
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-13",
      "connector_name": "test-connector-name-service-13",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-service-13",
      "state": 1,
      "output": "test-output-service-13",
      "client": "test-another-client-service-13"
    }
    """
    When I wait the end of 3 events processing
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "All: 0; Alarms: 0; Acknowledged: 0; NotAcknowledged: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "service.service",
                "m": "All: 0; Alarms: 0; Acknowledged: 0; NotAcknowledged: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given entity service and updated resource entity by api should update service alarm on infos update
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-14-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "infos.manager",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-manager-service-14"
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
      "connector": "test-connector-service-14",
      "connector_name": "test-connector-name-service-14",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-14",
      "resource": "test-resource-service-14",
      "state": 3,
      "output": "test-output-service-14"
    }
    """
    When I wait the end of event processing
    When I do PUT /api/v4/entitybasics?_id=test-resource-service-14/test-component-service-14:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "infos": [
        {
          "description": "Manager",
          "name": "manager",
          "value": "test-manager-service-14"
        }
      ],
      "impact": [
        "test-component-service-14"
      ],
      "depends": [
        "test-connector-service-14/test-connector-name-service-14"
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: should be able to retrieve dependencies and impacts
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-tree-1",
      "name": "test-service-tree-1",
      "output_template": "Test-service-tree-1",
      "category": "test-category-service-weather",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-nested-service-1",
                "test-nested-service-2",
                "test-nested-service-3"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-nested-service-1",
      "name": "test-nested-service-1",
      "output_template": "test-nested-service-1",
      "category": "test-category-service-weather",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-nested-service-resource-1",
                "test-nested-service-resource-2"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-nested-service-2",
      "name": "test-nested-service-2",
      "output_template": "test-nested-service-2",
      "category": "test-category-service-weather",
      "enabled": true,
      "impact_level": 2,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-nested-service-resource-3",
                "test-nested-service-resource-4"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-nested-service-3",
      "name": "test-nested-service-3",
      "output_template": "test-nested-service-3",
      "category": "test-category-service-weather",
      "enabled": true,
      "impact_level": 3,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-nested-service-resource-5",
                "test-nested-service-resource-6"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-nested-service-connector",
      "connector_name": "test-nested-service-connector-name",
      "component": "test-nested-service-component",
      "resource": "test-nested-service-resource-1",
      "source_type": "resource",
      "state": 1
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-nested-service-connector",
      "connector_name": "test-nested-service-connector-name",
      "component": "test-nested-service-component",
      "resource": "test-nested-service-resource-2",
      "source_type": "resource",
      "state": 2
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-nested-service-connector",
      "connector_name": "test-nested-service-connector-name",
      "component": "test-nested-service-component",
      "resource": "test-nested-service-resource-3",
      "source_type": "resource",
      "state": 3
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-nested-service-connector",
      "connector_name": "test-nested-service-connector-name",
      "component": "test-nested-service-component",
      "resource": "test-nested-service-resource-4",
      "source_type": "resource",
      "state": 0
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-nested-service-connector",
      "connector_name": "test-nested-service-connector-name",
      "component": "test-nested-service-component",
      "resource": "test-nested-service-resource-5",
      "source_type": "resource",
      "state": 0
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-nested-service-connector",
      "connector_name": "test-nested-service-connector-name",
      "component": "test-nested-service-component",
      "resource": "test-nested-service-resource-6",
      "source_type": "resource",
      "state": 0
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entityservice-dependencies?_id=test-service-tree-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "impact_state": 6,
          "has_dependencies": true,
          "alarm": {
            "d": "test-nested-service-2",
            "v": {
              "state": {
                "val": 3
              }
            }
          },
          "entity": {
            "_id": "test-nested-service-2",
            "impact_level": 2
          }
        },
        {
          "impact_state": 2,
          "has_dependencies": true,
          "alarm": {
            "d": "test-nested-service-1",
            "v": {
              "state": {
                "val": 2
              }
            }
          },
          "entity": {
            "_id": "test-nested-service-1",
            "impact_level": 1
          }
        },
        {
          "impact_state": 0,
          "has_dependencies": true,
          "alarm": null,
          "entity": {
            "_id": "test-nested-service-3",
            "impact_level": 3
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
    When I do GET /api/v4/entityservice-dependencies?_id=test-nested-service-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "impact_state": 2,
          "has_dependencies": false,
          "alarm": {
            "d": "test-nested-service-resource-2/test-nested-service-component",
            "v": {
              "state": {
                "val": 2
              }
            }
          },
          "entity": {
            "_id": "test-nested-service-resource-2/test-nested-service-component",
            "impact_level": 1
          }
        },
        {
          "impact_state": 1,
          "has_dependencies": false,
          "alarm": {
            "d": "test-nested-service-resource-1/test-nested-service-component",
            "v": {
              "state": {
                "val": 1
              }
            }
          },
          "entity": {
            "_id": "test-nested-service-resource-1/test-nested-service-component",
            "impact_level": 1
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
    When I do GET /api/v4/entityservice-dependencies?_id=test-nested-service-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "impact_state": 3,
          "has_dependencies": false,
          "alarm": {
            "d": "test-nested-service-resource-3/test-nested-service-component",
            "v": {
              "state": {
                "val": 3
              }
            }
          },
          "entity": {
            "_id": "test-nested-service-resource-3/test-nested-service-component",
            "impact_level": 1
          }
        },
        {
          "impact_state": 0,
          "has_dependencies": false,
          "alarm": null,
          "entity": {
            "_id": "test-nested-service-resource-4/test-nested-service-component",
            "impact_level": 1
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
    When I do GET /api/v4/entityservice-dependencies?_id=test-nested-service-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "impact_state": 0,
          "has_dependencies": false,
          "alarm": null,
          "entity": {
            "_id": "test-nested-service-resource-5/test-nested-service-component",
            "impact_level": 1
          }
        },
        {
          "impact_state": 0,
          "has_dependencies": false,
          "alarm": null,
          "entity": {
            "_id": "test-nested-service-resource-6/test-nested-service-component",
            "impact_level": 1
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
    When I do GET /api/v4/entityservice-dependencies?_id=test-nested-service-resource-1/test-nested-service-component
    Then the response code should be 404
    When I do GET /api/v4/entityservice-impacts?_id=test-service-tree-1
    Then the response code should be 200
    Then the response body should be:
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
    When I do GET /api/v4/entityservice-impacts?_id=test-nested-service-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "impact_state": 3,
          "has_impacts": false,
          "alarm": {
            "d": "test-service-tree-1",
            "v": {
              "state": {
                "val": 3
              }
            }
          },
          "entity": {
            "_id": "test-service-tree-1",
            "impact_level": 1
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
    When I do GET /api/v4/entityservice-impacts?_id=test-nested-service-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "impact_state": 3,
          "has_impacts": false,
          "alarm": {
            "d": "test-service-tree-1",
            "v": {
              "state": {
                "val": 3
              }
            }
          },
          "entity": {
            "_id": "test-service-tree-1",
            "impact_level": 1
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
    When I do GET /api/v4/entityservice-impacts?_id=test-nested-service-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "impact_state": 3,
          "has_impacts": false,
          "alarm": {
            "d": "test-service-tree-1",
            "v": {
              "state": {
                "val": 3
              }
            }
          },
          "entity": {
            "_id": "test-service-tree-1",
            "impact_level": 1
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
    When I do GET /api/v4/entityservice-impacts?_id=test-nested-service-resource-1/test-nested-service-component
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "impact_state": 2,
          "has_impacts": true,
          "alarm": {
            "d": "test-nested-service-1",
            "v": {
              "state": {
                "val": 2
              }
            }
          },
          "entity": {
            "_id": "test-nested-service-1",
            "impact_level": 1
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
    When I do GET /api/v4/entityservice-impacts?_id=test-nested-service-resource-3/test-nested-service-component
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "impact_state": 6,
          "has_impacts": true,
          "alarm": {
            "d": "test-nested-service-2",
            "v": {
              "state": {
                "val": 3
              }
            }
          },
          "entity": {
            "_id": "test-nested-service-2",
            "impact_level": 2
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
    When I do GET /api/v4/entityservice-impacts?_id=test-nested-service-resource-5/test-nested-service-component
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "impact_state": 0,
          "has_impacts": true,
          "alarm": null,
          "entity": {
            "_id": "test-nested-service-3",
            "impact_level": 3
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

  Scenario: given deleted entity service should delete service alarm
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-15-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-15"
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
      "connector": "test-connector-service-15",
      "connector_name": "test-connector-name-service-15",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-15",
      "resource": "test-resource-service-15",
      "state": 3,
      "output": "test-output-service-15"
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
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-15",
      "connector_name": "test-connector-name-service-15",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-15",
      "resource": "test-resource-service-15",
      "state": 2,
      "output": "test-output-service-15"
    }
    """
    When I wait the end of event processing
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

  Scenario: given disabled entity service should not update service alarm
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-16-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-16"
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
      "connector": "test-connector-service-16",
      "connector_name": "test-connector-name-service-16",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-16",
      "resource": "test-resource-service-16",
      "state": 3,
      "output": "test-output-service-16"
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
      "name": "test-entityservice-service-16-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": false,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-16"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-16",
      "connector_name": "test-connector-name-service-16",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-16",
      "resource": "test-resource-service-16",
      "state": 2,
      "output": "test-output-service-16"
    }
    """
    When I wait the end of event processing
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

  Scenario: given enabled entity service should create new service alarm
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-17-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-17"
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
      "connector": "test-connector-service-17",
      "connector_name": "test-connector-name-service-17",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-17",
      "resource": "test-resource-service-17",
      "state": 3,
      "output": "test-output-service-17"
    }
    """
    When I wait the end of 2 events processing
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-service-17-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": false,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-17"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    When I save response disableTimestamp={{ now }}
    Then the response code should be 200
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-17",
      "connector_name": "test-connector-name-service-17",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-17",
      "resource": "test-resource-service-17",
      "state": 1,
      "output": "test-output-service-17"
    }
    """
    When I wait the end of event processing
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-service-17-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-17"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 0; StateMajor: 0; StateMinor: 1; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given deleted entity service should update impact service alarm
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-18-name-1",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-18"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response dependServiceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-18-name-2",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-service-18-name-1"
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
      "connector": "test-connector-service-18",
      "connector_name": "test-connector-name-service-18",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-18",
      "resource": "test-resource-service-18",
      "state": 3,
      "output": "test-output-service-18"
    }
    """
    When I wait the end of 3 events processing
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "All: 0; Alarms: 0; Acknowledged: 0; NotAcknowledged: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "service.service",
                "m": "All: 0; Alarms: 0; Acknowledged: 0; NotAcknowledged: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given disabled entity service should update impact service alarm
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-19-name-1",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-19"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response dependServiceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-19-name-2",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-service-19-name-1"
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
      "connector": "test-connector-service-19",
      "connector_name": "test-connector-name-service-19",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-19",
      "resource": "test-resource-service-19",
      "state": 3,
      "output": "test-output-service-19"
    }
    """
    When I wait the end of 3 events processing
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
      "name": "test-entityservice-service-19-name-1",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": false,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-19"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "All: 0; Alarms: 0; Acknowledged: 0; NotAcknowledged: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "service.service",
                "m": "All: 0; Alarms: 0; Acknowledged: 0; NotAcknowledged: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given disabled entity should update service alarm
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-20-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-20"
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
      "connector": "test-connector-service-20",
      "connector_name": "test-connector-name-service-20",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-20",
      "resource": "test-resource-service-20",
      "state": 3,
      "output": "test-output-service-20"
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-service-20/test-component-service-20:
    """json
    {
      "impact_level": 1,
      "enabled": false,
      "infos": [],
      "impact": [
        "test-component-service-20"
      ],
      "depends": [
        "test-connector-service-20/test-connector-name-service-20"
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
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
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              },
              {
                "_t": "statedec",
                "a": "service.service",
                "m": "All: 0; Alarms: 0; Acknowledged: 0; NotAcknowledged: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "service.service",
                "m": "All: 0; Alarms: 0; Acknowledged: 0; NotAcknowledged: 0; StateCritical: 0; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
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

  Scenario: given new entity service shouldn't count double ack
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-service-21",
      "connector_name": "test-connector-name-service-21",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-service-21",
      "resource": "test-resource-service-21",
      "state": 3,
      "output": "test-output-service-21"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-21-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-21"
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
            "output": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];"
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
      "connector": "test-connector-service-21",
      "connector_name": "test-connector-name-service-21",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-service-21",
      "resource": "test-resource-service-21",
      "output": "test-output-service-21"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-21",
      "connector_name": "test-connector-name-service-21",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-service-21",
      "resource": "test-resource-service-21",
      "output": "test-output-service-21"
    }
    """
    When I wait the end of event processing
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
            "output": "All: 1; Alarms: 1; Acknowledged: 1; NotAcknowledged: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];"
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
