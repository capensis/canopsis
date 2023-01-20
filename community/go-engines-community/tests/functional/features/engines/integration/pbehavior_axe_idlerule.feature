Feature: update alarm on idle rule
  I need to be able to update alarm on idle rule

  Scenario: given pbehavior idle rule should update alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-pbehavior-axe-idlerule-1-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 50,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-axe-idlerule-1"
            }
          }
        ]
      ],
      "operation": {
        "type": "pbehavior",
        "parameters": {
          "name": "test-pbehavior-pbehavior-axe-idlerule-1",
          "start_on_trigger": true,
          "duration": {
            "value": 1,
            "unit": "h"
          },
          "type": "test-maintenance-type-to-engine",
          "reason": "test-reason-to-engine"
        }
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-pbehavior-axe-idlerule-1",
      "connector_name" : "test-connector-name-pbehavior-axe-idlerule-1",
      "source_type" : "resource",
      "component" :  "test-component-pbehavior-axe-idlerule-1",
      "resource" : "test-resource-pbehavior-axe-idlerule-1",
      "state" : 2,
      "output" : "test-output-pbehavior-axe-idlerule-1"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities/pbehaviors?_id=test-resource-pbehavior-axe-idlerule-1/test-component-pbehavior-axe-idlerule-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "name": "test-pbehavior-pbehavior-axe-idlerule-1"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-axe-idlerule-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-pbehavior-axe-idlerule-1",
            "connector": "test-connector-pbehavior-axe-idlerule-1",
            "connector_name": "test-connector-name-pbehavior-axe-idlerule-1",
            "resource": "test-resource-pbehavior-axe-idlerule-1",
            "pbehavior_info": {
              "canonical_type": "maintenance",
              "name": "test-pbehavior-pbehavior-axe-idlerule-1",
              "reason": "test-reason-to-engine",
              "reason_name": "Test Engine",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            },
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "pbhenter",
                "a": "system",
                "m": "Pbehavior test-pbehavior-pbehavior-axe-idlerule-1. Type: Engine maintenance. Reason: Test Engine."
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
    When I do GET /api/v4/entities/pbehaviors?_id=test-resource-pbehavior-axe-idlerule-1/test-component-pbehavior-axe-idlerule-1 until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-pbehavior-pbehavior-axe-idlerule-1"
      }
    ]
    """

  Scenario: given alarm idle rule should disable it on pbehavior
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-pbehavior-axe-idlerule-2",
      "connector_name" : "test-connector-name-pbehavior-axe-idlerule-2",
      "source_type" : "resource",
      "component" :  "test-component-pbehavior-axe-idlerule-2",
      "resource" : "test-resource-pbehavior-axe-idlerule-2",
      "state" : 0,
      "output" : "test-output-pbehavior-axe-idlerule-2"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-pbehavior-axe-idlerule-2-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 50,
      "duration": {
        "value": 1,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-axe-idlerule-2"
            }
          }
        ]
      ],
      "operation": {
        "type": "ack",
        "parameters": {
          "output": "test-pbehavior-pbehavior-axe-idlerule-2-output"
        }
      },
      "disable_during_periods": ["maintenance"]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "name": "test-pbehavior-pbehavior-axe-idlerule-2-name",
      "enabled": true,
      "tstart": {{ now }},
      "tstop": {{ nowAdd "7s" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-axe-idlerule-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-pbehavior-axe-idlerule-2",
      "connector_name" : "test-connector-name-pbehavior-axe-idlerule-2",
      "source_type" : "resource",
      "component" :  "test-component-pbehavior-axe-idlerule-2",
      "resource" : "test-resource-pbehavior-axe-idlerule-2",
      "state" : 2,
      "output" : "test-output-pbehavior-axe-idlerule-2"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-axe-idlerule-2 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-pbehavior-axe-idlerule-2",
            "connector": "test-connector-pbehavior-axe-idlerule-2",
            "connector_name": "test-connector-name-pbehavior-axe-idlerule-2",
            "resource": "test-resource-pbehavior-axe-idlerule-2",
            "pbehavior_info": {
              "canonical_type": "maintenance",
              "name": "test-pbehavior-pbehavior-axe-idlerule-2-name",
              "reason": "test-reason-to-engine",
              "reason_name": "Test Engine",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            },
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "pbhenter",
                "a": "system",
                "m": "Pbehavior test-pbehavior-pbehavior-axe-idlerule-2-name. Type: Engine maintenance. Reason: Test Engine."
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-axe-idlerule-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-pbehavior-axe-idlerule-2",
            "connector": "test-connector-pbehavior-axe-idlerule-2",
            "connector_name": "test-connector-name-pbehavior-axe-idlerule-2",
            "resource": "test-resource-pbehavior-axe-idlerule-2",
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "pbhenter",
                "a": "system",
                "m": "Pbehavior test-pbehavior-pbehavior-axe-idlerule-2-name. Type: Engine maintenance. Reason: Test Engine."
              },
              {
                "_t": "pbhleave",
                "a": "system",
                "m": "Pbehavior test-pbehavior-pbehavior-axe-idlerule-2-name. Type: Engine maintenance. Reason: Test Engine."
              },
              {
                "_t": "ack",
                "a": "system",
                "m": "test-pbehavior-pbehavior-axe-idlerule-2-output"
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

  Scenario: given entity idle rule should disable it on pbehavior and create alarm after pbehavior
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-pbehavior-axe-idlerule-3",
      "connector_name" : "test-connector-name-pbehavior-axe-idlerule-3",
      "source_type" : "resource",
      "component" :  "test-component-pbehavior-axe-idlerule-3",
      "resource" : "test-resource-pbehavior-axe-idlerule-3",
      "state" : 0,
      "output" : "test-output-pbehavior-axe-idlerule-3"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "name": "test-pbehavior-pbehavior-axe-idlerule-3-name",
      "enabled": true,
      "tstart": {{ now }},
      "tstop": {{ nowAdd "7s" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-axe-idlerule-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-pbehavior-axe-idlerule-3-name",
      "type": "entity",
      "enabled": true,
      "priority": 51,
      "duration": {
        "value": 1,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-axe-idlerule-3"
            }
          }
        ]
      ],
      "disable_during_periods": ["maintenance"]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-axe-idlerule-3 until response code is 200 and body contains:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-axe-idlerule-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-pbehavior-axe-idlerule-3",
            "connector": "test-connector-pbehavior-axe-idlerule-3",
            "connector_name": "test-connector-name-pbehavior-axe-idlerule-3",
            "resource": "test-resource-pbehavior-axe-idlerule-3",
            "state": {
              "val": 3
            },
            "status": {
              "val": 5
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
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
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

  Scenario: given entity idle rule should disable it on pbehavior and update alarm after pbehavior
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-pbehavior-axe-idlerule-4",
      "connector_name" : "test-connector-name-pbehavior-axe-idlerule-4",
      "source_type" : "resource",
      "component" :  "test-component-pbehavior-axe-idlerule-4",
      "resource" : "test-resource-pbehavior-axe-idlerule-4",
      "state" : 3,
      "output" : "test-output-pbehavior-axe-idlerule-4"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "name": "test-pbehavior-pbehavior-axe-idlerule-4-name",
      "enabled": true,
      "tstart": {{ now }},
      "tstop": {{ nowAdd "7s" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-axe-idlerule-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-pbehavior-axe-idlerule-4-name",
      "type": "entity",
      "enabled": true,
      "priority": 52,
      "duration": {
        "value": 1,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-axe-idlerule-4"
            }
          }
        ]
      ],
      "disable_during_periods": ["maintenance"]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-axe-idlerule-4 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-pbehavior-axe-idlerule-4",
            "connector": "test-connector-pbehavior-axe-idlerule-4",
            "connector_name": "test-connector-name-pbehavior-axe-idlerule-4",
            "resource": "test-resource-pbehavior-axe-idlerule-4",
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
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "pbhenter"
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-axe-idlerule-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-pbehavior-axe-idlerule-4",
            "connector": "test-connector-pbehavior-axe-idlerule-4",
            "connector_name": "test-connector-name-pbehavior-axe-idlerule-4",
            "resource": "test-resource-pbehavior-axe-idlerule-4",
            "state": {
              "val": 3
            },
            "status": {
              "val": 5
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
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "pbhenter"
              },
              {
                "_t": "pbhleave"
              },
              {
                "_t": "statusinc",
                "val": 5
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
