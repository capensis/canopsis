Feature: update alarm on pbehavior
  I need to be able to create pbehavior for alarm

  Scenario: given pbehavior should create alarm with pbehavior info
    Given I am admin
    When I send an event:
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
    When I wait the end of event processing
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
    When I wait the end of event processing
    When I send an event:
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
    When I wait the end of event processing
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
              "message": "Second comment"
            },
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

  Scenario: given pbehavior and alarm should update alarm pbehavior info
    Given I am admin
    When I send an event:
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
    When I wait the end of event processing
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
    When I wait the end of event processing
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

  Scenario: given pbehavior should update last alarm date of pbehavior
    Given I am admin
    When I send an event:
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
    When I wait the end of event processing
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
    When I wait the end of event processing
    When I do GET /api/v4/pbehaviors/{{ .pbehaviorID }}
    Then the response code should be 200
    Then the response key "last_alarm_date" should not be "null"

  Scenario: given pbehavior and entity without alarm should update last alarm date of pbehavior
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-5",
      "connector_name": "test-connector-name-pbehavior-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-5",
      "resource": "test-resource-pbehavior-5",
      "state": 0,
      "output": "test-output-pbehavior-5"
    }
    """
    When I wait the end of event processing
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
    Then the response body should contain:
    """json
    {
      "last_alarm_date": null
    }
    """
    When I save response pbehaviorID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do GET /api/v4/pbehaviors/{{ .pbehaviorID }}
    Then the response code should be 200
    Then the response key "last_alarm_date" should not be "null"

  Scenario: given deleted pbehavior should delete alarm with pbehavior info
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-6",
      "connector_name": "test-connector-name-pbehavior-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-6",
      "resource": "test-resource-pbehavior-6",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
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
              "value": "test-resource-pbehavior-6"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do DELETE /api/v4/pbehaviors/{{ .lastResponse._id }}
    Then the response code should be 204
    When I wait the end of event processing
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

  Scenario: given updated pbehavior entity pattern should delete alarm with pbehavior info
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-7",
      "connector_name": "test-connector-name-pbehavior-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-7",
      "resource": "test-resource-pbehavior-7-1",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-7",
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
              "value": "test-resource-pbehavior-7-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do PUT /api/v4/pbehaviors/{{ .lastResponse._id }}:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-7",
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
              "value": "test-resource-pbehavior-7-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-7
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
                "m": "Pbehavior test-pbehavior-7. Type: Engine maintenance. Reason: Test Engine."
              },
              {
                "_t": "pbhleave",
                "a": "system",
                "user_id": "",
                "m": "Pbehavior test-pbehavior-7. Type: Engine maintenance. Reason: Test Engine."
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

  Scenario: given pbehavior and alarm should update alarm pbehavior info on periodical
    Given I am admin
    When I send an event:
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
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-8",
      "tstart": {{ nowAdd "2s" }},
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
    When I wait the end of event processing
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

  Scenario: given pbehavior should create alarm with pbehavior info
    Given I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-9",
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
              "value": "test-resource-pbehavior-9"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-9",
      "connector_name": "test-connector-name-pbehavior-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-9",
      "resource": "test-resource-pbehavior-9",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test-pbehavior-9",
              "reason": "test-reason-to-engine",
              "reason_name": "Test Engine",
              "canonical_type": "maintenance",
              "icon_name": "build",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            },
            "connector": "test-connector-pbehavior-9",
            "connector_name": "test-connector-name-pbehavior-9",
            "component": "test-component-pbehavior-9",
            "resource": "test-resource-pbehavior-9"
          },
          "pbehavior": {
            "name": "test-pbehavior-9",
            "last_comment": null,
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
                "m": "Pbehavior test-pbehavior-9. Type: Engine maintenance. Reason: Test Engine."
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

  Scenario: given pbehavior with pause type and without stop should create alarm with pbehavior info
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-10",
      "connector_name": "test-connector-name-pbehavior-10",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-10",
      "resource": "test-resource-pbehavior-10",
      "state": 0,
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-10",
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
              "value": "test-resource-pbehavior-10"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-pbehavior-10
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-10",
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

  Scenario: given pbehavior should create entity and alarm with pbehavior info
    Given I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-11",
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
              "value": "test-resource-pbehavior-11"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-11 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-11",
          "is_active_status": true
        }
      ]
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-11",
      "connector_name": "test-connector-name-pbehavior-11",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-11",
      "resource": "test-resource-pbehavior-11",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-11
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test-pbehavior-11",
              "reason": "test-reason-to-engine",
              "reason_name": "Test Engine",
              "canonical_type": "maintenance",
              "icon_name": "build",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            },
            "connector": "test-connector-pbehavior-11",
            "connector_name": "test-connector-name-pbehavior-11",
            "component": "test-component-pbehavior-11",
            "resource": "test-resource-pbehavior-11"
          },
          "pbehavior": {
            "name": "test-pbehavior-11",
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
                "m": "Pbehavior test-pbehavior-11. Type: Engine maintenance. Reason: Test Engine."
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
    When I do GET /api/v4/entities?search=test-resource-pbehavior-11
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-11",
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

  Scenario: given pbehavior should create entity with pbehavior info
    Given I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-12",
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
              "value": "test-resource-pbehavior-12"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-12 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-12",
          "is_active_status": true
        }
      ]
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-12",
      "connector_name": "test-connector-name-pbehavior-12",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-12",
      "resource": "test-resource-pbehavior-12",
      "state": 0,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-12
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
    When I do GET /api/v4/entities?search=test-resource-pbehavior-12
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-12",
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

  Scenario: given updated corporate entity pattern should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-13",
      "connector_name": "test-connector-name-pbehavior-13",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-13",
      "resource": "test-resource-pbehavior-13",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/patterns:
    """json
    {
      "title": "test-pattern-pbehavior-13",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-13"
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
      "name": "test-pbehavior-13",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "corporate_entity_pattern": "{{ .patternID }}"
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-pbehavior-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-13",
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
      "title": "test-pattern-pbehavior-13",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-13-not-exist"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-pbehavior-13
    Then the response code should be 200
    Then the response key "data.0.pheavior_info" should not exist
    When I do PUT /api/v4/patterns/{{ .patternID }}:
    """json
    {
      "title": "test-pattern-pbehavior-13",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-13"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-pbehavior-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-13",
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
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test-pbehavior-13",
              "reason": "test-reason-to-engine",
              "reason_name": "Test Engine",
              "canonical_type": "maintenance",
              "icon_name": "build",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            },
            "connector": "test-connector-pbehavior-13",
            "connector_name": "test-connector-name-pbehavior-13",
            "component": "test-component-pbehavior-13",
            "resource": "test-resource-pbehavior-13"
          },
          "pbehavior": {
            "name": "test-pbehavior-13",
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
                "m": "Pbehavior test-pbehavior-13. Type: Engine maintenance. Reason: Test Engine."
              },
              {
                "_t": "pbhleave",
                "a": "system",
                "user_id": "",
                "m": "Pbehavior test-pbehavior-13. Type: Engine maintenance. Reason: Test Engine."
              },
              {
                "_t": "pbhenter",
                "a": "system",
                "user_id": "",
                "m": "Pbehavior test-pbehavior-13. Type: Engine maintenance. Reason: Test Engine."
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

  Scenario: given pbehavior with old mongo query should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-14",
      "connector_name": "test-connector-name-pbehavior-14",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-14",
      "resource": "test-resource-pbehavior-14",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do PUT /api/v4/pbehaviors/test-pbehavior-14:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-14",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-14
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test-pbehavior-14",
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
                "m": "Pbehavior test-pbehavior-14. Type: Engine maintenance. Reason: Test Engine."
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
