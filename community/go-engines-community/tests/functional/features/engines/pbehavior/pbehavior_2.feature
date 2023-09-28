Feature: update alarm on pbehavior
  I need to be able to create pbehavior for alarm

  @concurrent
  Scenario: given pbehavior should create entity and alarm with pbehavior info
    Given I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-second-1",
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
              "value": "test-resource-pbehavior-second-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-second-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-second-1",
          "is_active_status": true
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-second-1",
      "connector_name": "test-connector-name-pbehavior-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-second-1",
      "resource": "test-resource-pbehavior-second-1",
      "state": 1,
      "output": "test-output-pbehavior-second-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-second-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test-pbehavior-second-1",
              "reason": "test-reason-to-engine",
              "reason_name": "Test Engine",
              "canonical_type": "maintenance",
              "icon_name": "build",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            },
            "connector": "test-connector-pbehavior-second-1",
            "connector_name": "test-connector-name-pbehavior-second-1",
            "component": "test-component-pbehavior-second-1",
            "resource": "test-resource-pbehavior-second-1"
          },
          "pbehavior": {
            "name": "test-pbehavior-second-1",
            "author": {
              "_id": "root",
              "name": "root"
            },
            "reason": {
              "_id": "test-reason-to-engine",
              "name": "Test Engine",
              "description": "Test Engine"
            },
            "type": {
              "_id": "test-maintenance-type-to-engine",
              "icon_name": "build",
              "name": "Engine maintenance",
              "type": "maintenance"
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
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "pbhenter",
                "a": "system",
                "user_id": "",
                "m": "Pbehavior test-pbehavior-second-1. Type: Engine maintenance. Reason: Test Engine."
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
    When I do GET /api/v4/entities?search=test-resource-pbehavior-second-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-second-1",
            "reason": "test-reason-to-engine",
            "reason_name": "Test Engine",
            "canonical_type": "maintenance",
            "icon_name": "build",
            "type": "test-maintenance-type-to-engine",
            "type_name": "Engine maintenance"
          }
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """

  @concurrent
  Scenario: given pbehavior should create entity with pbehavior info
    Given I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-second-2",
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
              "value": "test-resource-pbehavior-second-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-second-2 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-second-2",
          "is_active_status": true
        }
      ]
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-second-2",
      "connector_name": "test-connector-name-pbehavior-second-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-second-2",
      "resource": "test-resource-pbehavior-second-2",
      "state": 0,
      "output": "test-output-pbehavior-second-2"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-pbehavior-second-2",
        "connector_name": "test-connector-name-pbehavior-second-2",
        "component": "test-component-pbehavior-second-2",
        "resource": "test-resource-pbehavior-second-2",
        "source_type": "resource"
      },
      {
        "event_type": "pbhenter",
        "connector": "test-connector-pbehavior-second-2",
        "connector_name": "test-connector-name-pbehavior-second-2",
        "component": "test-component-pbehavior-second-2",
        "resource": "test-resource-pbehavior-second-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-second-2
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
    When I do GET /api/v4/entities?search=test-resource-pbehavior-second-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-second-2",
            "reason": "test-reason-to-engine",
            "reason_name": "Test Engine",
            "canonical_type": "maintenance",
            "icon_name": "build",
            "type": "test-maintenance-type-to-engine",
            "type_name": "Engine maintenance"
          }
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """

  @concurrent
  Scenario: given updated corporate entity pattern should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-pbehavior-second-3",
      "connector": "test-connector-pbehavior-second-3",
      "connector_name": "test-connector-name-pbehavior-second-3",
      "component": "test-component-pbehavior-second-3",
      "resource": "test-resource-pbehavior-second-3",
      "source_type": "resource"
    }
    """
    When I do POST /api/v4/patterns:
    """json
    {
      "title": "test-pattern-pbehavior-second-3",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-second-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response patternID={{ .lastResponse._id }}
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-second-3",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "corporate_entity_pattern": "{{ .patternID }}"
    }
    """
    Then the response code should be 201
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "pbhenter",
      "connector": "test-connector-pbehavior-second-3",
      "connector_name": "test-connector-name-pbehavior-second-3",
      "component": "test-component-pbehavior-second-3",
      "resource": "test-resource-pbehavior-second-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/entities?search=test-resource-pbehavior-second-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-second-3",
            "reason": "test-reason-to-engine",
            "reason_name": "Test Engine",
            "canonical_type": "maintenance",
            "icon_name": "build",
            "type": "test-maintenance-type-to-engine",
            "type_name": "Engine maintenance"
          }
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """
    When I do PUT /api/v4/patterns/{{ .patternID }}:
    """json
    {
      "title": "test-pattern-pbehavior-second-3",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-second-3-not-exist"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "pbhleave",
      "connector": "test-connector-pbehavior-second-3",
      "connector_name": "test-connector-name-pbehavior-second-3",
      "component": "test-component-pbehavior-second-3",
      "resource": "test-resource-pbehavior-second-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/entities?search=test-resource-pbehavior-second-3
    Then the response code should be 200
    Then the response key "data.0.pheavior_info" should not exist
    When I do PUT /api/v4/patterns/{{ .patternID }}:
    """json
    {
      "title": "test-pattern-pbehavior-second-3",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-second-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "pbhenter",
      "connector": "test-connector-pbehavior-second-3",
      "connector_name": "test-connector-name-pbehavior-second-3",
      "component": "test-component-pbehavior-second-3",
      "resource": "test-resource-pbehavior-second-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/entities?search=test-resource-pbehavior-second-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-second-3",
            "reason": "test-reason-to-engine",
            "reason_name": "Test Engine",
            "canonical_type": "maintenance",
            "icon_name": "build",
            "type": "test-maintenance-type-to-engine",
            "type_name": "Engine maintenance"
          }
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-second-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test-pbehavior-second-3",
              "reason": "test-reason-to-engine",
              "reason_name": "Test Engine",
              "canonical_type": "maintenance",
              "icon_name": "build",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            },
            "connector": "test-connector-pbehavior-second-3",
            "connector_name": "test-connector-name-pbehavior-second-3",
            "component": "test-component-pbehavior-second-3",
            "resource": "test-resource-pbehavior-second-3"
          },
          "pbehavior": {
            "name": "test-pbehavior-second-3",
            "author": {
              "_id": "root",
              "name": "root"
            },
            "reason": {
              "_id": "test-reason-to-engine",
              "name": "Test Engine",
              "description": "Test Engine"
            },
            "type": {
              "_id": "test-maintenance-type-to-engine",
              "icon_name": "build",
              "name": "Engine maintenance",
              "type": "maintenance"
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
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "pbhenter",
                "a": "system",
                "user_id": "",
                "m": "Pbehavior test-pbehavior-second-3. Type: Engine maintenance. Reason: Test Engine."
              },
              {
                "_t": "pbhleave",
                "a": "system",
                "user_id": "",
                "m": "Pbehavior test-pbehavior-second-3. Type: Engine maintenance. Reason: Test Engine."
              },
              {
                "_t": "pbhenter",
                "a": "system",
                "user_id": "",
                "m": "Pbehavior test-pbehavior-second-3. Type: Engine maintenance. Reason: Test Engine."
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 5
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given pbehavior with old mongo query should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-second-4",
      "connector_name": "test-connector-name-pbehavior-second-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-second-4",
      "resource": "test-resource-pbehavior-second-4",
      "state": 1,
      "output": "test-output-pbehavior-second-4"
    }
    """
    When I do PUT /api/v4/pbehaviors/test-pbehavior-second-4:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-second-4",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "pbhenter",
      "connector": "test-connector-pbehavior-second-4",
      "connector_name": "test-connector-name-pbehavior-second-4",
      "component": "test-component-pbehavior-second-4",
      "resource": "test-resource-pbehavior-second-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-second-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test-pbehavior-second-4",
              "reason": "test-reason-to-engine",
              "reason_name": "Test Engine",
              "canonical_type": "maintenance",
              "icon_name": "build",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
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
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "pbhenter",
                "a": "system",
                "user_id": "",
                "m": "Pbehavior test-pbehavior-second-4. Type: Engine maintenance. Reason: Test Engine."
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
  Scenario: given entity with pbehavior_info should set creation date in alarm's pbehavior info
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 0,
      "output": "noveo alarm",
      "connector": "test-connector-pbehavior-second-5",
      "connector_name": "test-connector-name-pbehavior-second-5",
      "component": "test-component-pbehavior-second-5",
      "resource": "test-resource-pbehavior-second-5",
      "source_type": "resource"
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-second-5",
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
              "value": "test-resource-pbehavior-second-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "pbhenter",
      "connector": "test-connector-pbehavior-second-5",
      "connector_name": "test-connector-name-pbehavior-second-5",
      "component": "test-component-pbehavior-second-5",
      "resource": "test-resource-pbehavior-second-5",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/entities?search=test-resource-pbehavior-second-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-second-5"
          }
        }
      ]
    }
    """
    When I save response entityPbhTs={{ (index .lastResponse.data 0).pbehavior_info.timestamp }}
    When I wait 2s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "noveo alarm",
      "connector": "test-connector-pbehavior-second-5",
      "connector_name": "test-connector-name-pbehavior-second-5",
      "component": "test-component-pbehavior-second-5",
      "resource": "test-resource-pbehavior-second-5",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-second-5
    Then the response code should be 200
    When I save response alarmCreationDate={{ (index .lastResponse.data 0).t }}
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "timestamp": {{ .alarmCreationDate }},
              "name": "test-pbehavior-second-5"
            }
          }
        }
      ]
    }
    """
    When I save response alarmPbhTs={{ (index .lastResponse.data 0).v.pbehavior_info.timestamp }}
    Then "entityPbhTs" < "alarmPbhTs"
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
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "pbhenter",
                "t": {{ .alarmCreationDate }}
              }
            ]
          }
        }
      }
    ]
    """
    When I do GET /api/v4/entities?search=test-resource-pbehavior-second-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "timestamp": {{ .entityPbhTs }},
            "name": "test-pbehavior-second-5"
          }
        }
      ]
    }
    """
