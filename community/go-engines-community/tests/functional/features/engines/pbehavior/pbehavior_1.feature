Feature: update alarm on pbehavior
  I need to be able to create pbehavior for alarm

  @concurrent
  Scenario: given pbehavior should create alarm with pbehavior info
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-1",
      "connector_name": "test-connector-name-pbehavior-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-1",
      "resource": "test-resource-pbehavior-1",
      "state": 0,
      "output": "noveo alarm"
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-1",
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
              "value": "test-resource-pbehavior-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response pbehaviorID={{ .lastResponse._id }}
    When I do POST /api/v4/pbehavior-comments:
    """json
    {
      "pbehavior": "{{ .pbehaviorID }}",
      "message": "First comment"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehavior-comments:
    """json
    {
      "pbehavior": "{{ .pbehaviorID }}",
      "message": "Second comment"
    }
    """
    Then the response code should be 201
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "pbhenter",
      "connector": "test-connector-pbehavior-1",
      "connector_name": "test-connector-name-pbehavior-1",
      "component": "test-component-pbehavior-1",
      "resource": "test-resource-pbehavior-1",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-1",
      "connector_name": "test-connector-name-pbehavior-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-1",
      "resource": "test-resource-pbehavior-1",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test-pbehavior-1",
              "reason": "test-reason-to-engine",
              "reason_name": "Test Engine",
              "canonical_type": "maintenance",
              "icon_name": "build",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            },
            "connector": "test-connector-pbehavior-1",
            "connector_name": "test-connector-name-pbehavior-1",
            "component": "test-component-pbehavior-1",
            "resource": "test-resource-pbehavior-1"
          },
          "pbehavior": {
            "name": "test-pbehavior-1",
            "last_comment": {
              "message": "Second comment",
              "author": {
                "_id": "root",
                "name": "root",
                "display_name": "root John Doe admin@canopsis.net"
              }
            },
            "author": {
              "_id": "root",
              "name": "root",
              "display_name": "root John Doe admin@canopsis.net"
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
                "m": "Pbehavior test-pbehavior-1. Type: Engine maintenance. Reason: Test Engine."
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
    When I do GET /api/v4/entities?search=test-resource-pbehavior-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-1",
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
  Scenario: given pbehavior and alarm should update alarm pbehavior info
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-2",
      "connector_name": "test-connector-name-pbehavior-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-2",
      "resource": "test-resource-pbehavior-2",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-2",
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
              "value": "test-resource-pbehavior-2"
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
      "event_type" : "pbhenter",
      "connector": "test-connector-pbehavior-2",
      "connector_name": "test-connector-name-pbehavior-2",
      "component": "test-component-pbehavior-2",
      "resource": "test-resource-pbehavior-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test-pbehavior-2",
              "reason": "test-reason-to-engine",
              "reason_name": "Test Engine",
              "canonical_type": "maintenance",
              "icon_name": "build",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            },
            "connector": "test-connector-pbehavior-2",
            "connector_name": "test-connector-name-pbehavior-2",
            "component": "test-component-pbehavior-2",
            "resource": "test-resource-pbehavior-2"
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
                "m": "Pbehavior test-pbehavior-2. Type: Engine maintenance. Reason: Test Engine."
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
    When I do GET /api/v4/entities?search=test-resource-pbehavior-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-2",
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
  Scenario: given pbehavior should update last alarm date of pbehavior
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-3",
      "connector_name": "test-connector-name-pbehavior-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-3",
      "resource": "test-resource-pbehavior-3",
      "state": 1,
      "output": "test-output-pbehavior-3"
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-3",
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
              "value": "test-resource-pbehavior-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "last_alarm_date": null
    }
    """
    When I save response pbehaviorID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "pbhenter",
      "connector": "test-connector-pbehavior-3",
      "connector_name": "test-connector-name-pbehavior-3",
      "component": "test-component-pbehavior-3",
      "resource": "test-resource-pbehavior-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/pbehaviors/{{ .pbehaviorID }} until response code is 200 and response key "last_alarm_date" is greater or equal than 1

  @concurrent
  Scenario: given pbehavior and entity without alarm should update last alarm date of pbehavior
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-4",
      "connector_name": "test-connector-name-pbehavior-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-4",
      "resource": "test-resource-pbehavior-4",
      "state": 0,
      "output": "test-output-pbehavior-4"
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-4",
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
              "value": "test-resource-pbehavior-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "last_alarm_date": null
    }
    """
    When I save response pbehaviorID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "pbhenter",
      "connector": "test-connector-pbehavior-4",
      "connector_name": "test-connector-name-pbehavior-4",
      "component": "test-component-pbehavior-4",
      "resource": "test-resource-pbehavior-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/pbehaviors/{{ .pbehaviorID }} until response code is 200 and response key "last_alarm_date" is greater or equal than 1

  @concurrent
  Scenario: given deleted pbehavior should delete alarm with pbehavior info
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-5",
      "connector_name": "test-connector-name-pbehavior-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-5",
      "resource": "test-resource-pbehavior-5",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-5",
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
              "value": "test-resource-pbehavior-5"
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
      "event_type" : "pbhenter",
      "connector": "test-connector-pbehavior-5",
      "connector_name": "test-connector-name-pbehavior-5",
      "component": "test-component-pbehavior-5",
      "resource": "test-resource-pbehavior-5",
      "source_type": "resource"
    }
    """
    When I do DELETE /api/v4/pbehaviors/{{ .lastResponse._id }}
    Then the response code should be 204
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "pbhleave",
      "connector": "test-connector-pbehavior-5",
      "connector_name": "test-connector-name-pbehavior-5",
      "component": "test-component-pbehavior-5",
      "resource": "test-resource-pbehavior-5",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-5
    Then the response code should be 200
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "pbhenter",
                "a": "system",
                "user_id": "",
                "m": "Pbehavior test-pbehavior-5. Type: Engine maintenance. Reason: Test Engine."
              },
              {
                "_t": "pbhleave",
                "a": "system",
                "user_id": "",
                "m": "Pbehavior test-pbehavior-5. Type: Engine maintenance. Reason: Test Engine."
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
  Scenario: given updated pbehavior entity pattern should delete alarm with pbehavior info
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-6",
      "connector_name": "test-connector-name-pbehavior-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-6",
      "resource": "test-resource-pbehavior-6-1",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-6",
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
              "value": "test-resource-pbehavior-6-1"
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
      "event_type" : "pbhenter",
      "connector": "test-connector-pbehavior-6",
      "connector_name": "test-connector-name-pbehavior-6",
      "component": "test-component-pbehavior-6",
      "resource": "test-resource-pbehavior-6-1",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/pbehaviors/{{ .lastResponse._id }}:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-6",
      "tstart": {{ .lastResponse.tstart }},
      "tstop": {{ .lastResponse.tstop }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-6-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "pbhleave",
      "connector": "test-connector-pbehavior-6",
      "connector_name": "test-connector-name-pbehavior-6",
      "component": "test-component-pbehavior-6",
      "resource": "test-resource-pbehavior-6-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-6
    Then the response code should be 200
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "pbhenter",
                "a": "system",
                "user_id": "",
                "m": "Pbehavior test-pbehavior-6. Type: Engine maintenance. Reason: Test Engine."
              },
              {
                "_t": "pbhleave",
                "a": "system",
                "user_id": "",
                "m": "Pbehavior test-pbehavior-6. Type: Engine maintenance. Reason: Test Engine."
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
  Scenario: given pbehavior and alarm should update alarm pbehavior info on periodical
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-7",
      "connector_name": "test-connector-name-pbehavior-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-7",
      "resource": "test-resource-pbehavior-7",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-7",
      "tstart": {{ nowAdd "3s" }},
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
              "value": "test-resource-pbehavior-7"
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
      "event_type" : "pbhenter",
      "connector": "test-connector-pbehavior-7",
      "connector_name": "test-connector-name-pbehavior-7",
      "component": "test-component-pbehavior-7",
      "resource": "test-resource-pbehavior-7",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test-pbehavior-7",
              "reason": "test-reason-to-engine",
              "reason_name": "Test Engine",
              "canonical_type": "maintenance",
              "icon_name": "build",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            },
            "connector": "test-connector-pbehavior-7",
            "connector_name": "test-connector-name-pbehavior-7",
            "component": "test-component-pbehavior-7",
            "resource": "test-resource-pbehavior-7"
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
                "m": "Pbehavior test-pbehavior-7. Type: Engine maintenance. Reason: Test Engine."
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
  Scenario: given pbehavior should create alarm with pbehavior info
    Given I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-8",
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
              "value": "test-resource-pbehavior-8"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-8",
      "connector_name": "test-connector-name-pbehavior-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-8",
      "resource": "test-resource-pbehavior-8",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test-pbehavior-8",
              "reason": "test-reason-to-engine",
              "reason_name": "Test Engine",
              "canonical_type": "maintenance",
              "icon_name": "build",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            },
            "connector": "test-connector-pbehavior-8",
            "connector_name": "test-connector-name-pbehavior-8",
            "component": "test-component-pbehavior-8",
            "resource": "test-resource-pbehavior-8"
          },
          "pbehavior": {
            "name": "test-pbehavior-8",
            "last_comment": null,
            "author": {
              "_id": "root",
              "name": "root",
              "display_name": "root John Doe admin@canopsis.net"
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
                "m": "Pbehavior test-pbehavior-8. Type: Engine maintenance. Reason: Test Engine."
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
    When I do GET /api/v4/entities?search=test-resource-pbehavior-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-8",
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
  Scenario: given pbehavior with pause type and without stop should create alarm with pbehavior info
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-9",
      "connector_name": "test-connector-name-pbehavior-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-9",
      "resource": "test-resource-pbehavior-9",
      "state": 0,
      "output": "noveo alarm"
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-9",
      "tstart": {{ now }},
      "tstop": null,
      "color": "#FFFFFF",
      "type": "test-pause-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-9"
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
      "event_type" : "pbhenter",
      "connector": "test-connector-pbehavior-9",
      "connector_name": "test-connector-name-pbehavior-9",
      "component": "test-component-pbehavior-9",
      "resource": "test-resource-pbehavior-9",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/entities?search=test-resource-pbehavior-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-9",
            "reason": "test-reason-to-engine",
            "reason_name": "Test Engine",
            "canonical_type": "pause",
            "icon_name": "pause",
            "type": "test-pause-type-to-engine",
            "type_name": "Engine pause"
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
