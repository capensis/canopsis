Feature: update meta alarm on action
  I need to be able to update meta alarm on action

  Scenario: given meta alarm and scenario should update meta alarm and update children
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-action-correlation-1
    When I send an event:
    """json
    {
      "connector": "test-connector-action-correlation-1",
      "connector_name": "test-connector-name-action-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-action-correlation-1",
      "resource": "test-resource-action-correlation-1",
      "state": 2,
      "output": "test-output-action-correlation-1"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID }}&active_columns[]=v.meta&with_consequences=true&correlation=true
    Then the response code should be 200
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-correlation-1-name",
      "enabled": true,
      "triggers": ["comment"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "{{ .metaAlarmResource }}"
                }
              }
            ]
          ],
          "type": "assocticket",
          "parameters": {
            "output": "test-output-action-correlation-1-{{ `{{ .Alarm.Value.Connector }}` }}",
            "ticket": "test-ticket-action-correlation-1"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "{{ .metaAlarmResource }}"
                }
              }
            ]
          ],
          "type": "ack",
          "parameters": {
            "output": "test-output-action-correlation-1-{{ `{{ .Alarm.Value.Connector }}` }}"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "comment",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-output-action-correlation-1"
    }
    """
    When I wait the end of 4 events processing
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "component": "test-component-action-correlation-1",
                  "connector": "test-connector-action-correlation-1",
                  "connector_name": "test-connector-name-action-correlation-1",
                  "resource": "test-resource-action-correlation-1"
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
        }
      }
    ]
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 0)._id }}",
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
                "_t": "comment"
              },
              {
                "_t": "assocticket",
                "a": "system",
                "m": "test-ticket-action-correlation-1"
              },
              {
                "_t": "ack",
                "a": "system",
                "m": "test-output-action-correlation-1-engine"
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
      },
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
                "_t": "metaalarmattach",
                "val": 0
              },
              {
                "_t": "comment"
              },
              {
                "_t": "assocticket",
                "a": "system",
                "m": "test-ticket-action-correlation-1"
              },
              {
                "_t": "ack",
                "a": "system",
                "m": "test-output-action-correlation-1-engine"
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

  Scenario: given meta alarm and scenario with webhook action should update meta alarm and update children
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-action-correlation-2
    When I send an event:
    """json
    {
      "connector": "test-connector-action-correlation-2",
      "connector_name": "test-connector-name-action-correlation-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-action-correlation-2",
      "resource": "test-resource-action-correlation-2-1",
      "state": 2,
      "output": "test-output-action-correlation-2"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-action-correlation-2",
      "connector_name": "test-connector-name-action-correlation-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-action-correlation-2",
      "resource": "test-resource-action-correlation-2-2",
      "state": 1,
      "output": "test-output-action-correlation-2"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID }}&active_columns[]=v.meta&with_consequences=true&correlation=true
    Then the response code should be 200
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaalarmDisplayName={{ (index .lastResponse.data 0).v.display_name }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-correlation-2-name",
      "enabled": true,
      "triggers": ["comment"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "{{ .metaAlarmResource }}"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "http://localhost:3000/webhook/ticket"
            },
            "declare_ticket": {
              "ticket_id": "ticket_id",
              "ticket_data": "ticket_data"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "comment",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-output-action-correlation-2"
    }
    """
    When I wait the end of 5 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-action-correlation-2",
      "connector_name": "test-connector-name-action-correlation-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-action-correlation-2",
      "resource": "test-resource-action-correlation-2-3",
      "state": 2,
      "output": "test-output-action-correlation-2"
    }
    """
    When I wait the end of 3 events processing
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "component": "test-component-action-correlation-2",
                  "connector": "test-connector-action-correlation-2",
                  "connector_name": "test-connector-name-action-correlation-2",
                  "resource": "test-resource-action-correlation-2-1",
                  "ticket": {
                    "_t": "declareticket",
                    "val": "testticket",
                    "data": {
                      "ticket_data": "testdata"
                    }
                  }
                }
              },
              {
                "v": {
                  "component": "test-component-action-correlation-2",
                  "connector": "test-connector-action-correlation-2",
                  "connector_name": "test-connector-name-action-correlation-2",
                  "resource": "test-resource-action-correlation-2-2",
                  "ticket": {
                    "_t": "declareticket",
                    "val": "testticket",
                    "data": {
                      "ticket_data": "testdata"
                    }
                  }
                }
              },
              {
                "v": {
                  "component": "test-component-action-correlation-2",
                  "connector": "test-connector-action-correlation-2",
                  "connector_name": "test-connector-name-action-correlation-2",
                  "resource": "test-resource-action-correlation-2-3",
                  "ticket": {
                    "_t": "declareticket",
                    "val": "testticket",
                    "data": {
                      "ticket_data": "testdata"
                    }
                  }
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
        }
      }
    ]
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 0)._id }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 1)._id }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 2)._id }}",
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
                "_t": "comment"
              },
              {
                "_t": "declareticket"
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
      },
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "val": 0
              },
              {
                "_t": "comment"
              },
              {
                "_t": "declareticket"
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
      },
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "val": 0
              },
              {
                "_t": "comment"
              },
              {
                "_t": "declareticket"
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
      },
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "val": 0
              },
              {
                "_t": "declareticket"
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

  Scenario: given scenario for meta alarm and scenario for child alarm should update child alarm in correct order
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-action-correlation-3
    When I send an event:
    """json
    {
      "connector": "test-connector-action-correlation-3",
      "connector_name": "test-connector-name-action-correlation-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-action-correlation-3",
      "resource": "test-resource-action-correlation-3",
      "state": 2,
      "output": "test-output-action-correlation-3"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID }}&active_columns[]=v.meta&with_consequences=true&correlation=true
    Then the response code should be 200
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaalarmDisplayName={{ (index .lastResponse.data 0).v.display_name }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-correlation-3-1-name",
      "enabled": true,
      "triggers": ["comment"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "{{ .metaAlarmResource }}"
                }
              }
            ]
          ],
          "type": "ack",
          "parameters": {
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-correlation-3-2-name",
      "enabled": true,
      "triggers": ["comment"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-correlation-3"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "duration": {
              "value": 1,
              "unit": "h"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    [
      {
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "source_type": "resource",
        "event_type": "comment",
        "component":  "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "output": "test-output-action-correlation-3-metaalarm"
      },
      {
        "connector": "test-connector-action-correlation-3",
        "connector_name": "test-connector-name-action-correlation-3",
        "source_type": "resource",
        "event_type": "comment",
        "component":  "test-component-action-correlation-3",
        "resource": "test-resource-action-correlation-3",
        "output": "test-output-action-correlation-3"
      }
    ]
    """
    When I wait the end of 4 events processing
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "component": "test-component-action-correlation-3",
                  "connector": "test-connector-action-correlation-3",
                  "connector_name": "test-connector-name-action-correlation-3",
                  "resource": "test-resource-action-correlation-3"
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
        }
      }
    ]
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 0)._id }}",
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
                "_t": "comment"
              },
              {
                "_t": "ack"
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
      },
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
                "_t": "metaalarmattach",
                "val": 0
              },
              {
                "_t": "comment",
                "m": "test-output-action-correlation-3"
              },
              {
                "_t": "snooze"
              },
              {
                "_t": "comment",
                "m": "test-output-action-correlation-3-metaalarm"
              },
              {
                "_t": "ack"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 7
            }
          }
        }
      }
    ]
    """

  Scenario: given meta alarm and scenario for child alarm should update meta alarm state
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-action-correlation-4
    When I send an event:
    """json
    {
      "connector": "test-connector-action-correlation-4",
      "connector_name": "test-connector-name-action-correlation-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-action-correlation-4",
      "resource": "test-resource-action-correlation-4",
      "state": 2,
      "output": "test-output-action-correlation-4"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID }}&active_columns[]=v.meta&with_consequences=true&correlation=true
    Then the response code should be 200
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-correlation-4-name",
      "enabled": true,
      "triggers": ["comment"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-correlation-4"
                }
              }
            ]
          ],
          "type": "changestate",
          "parameters": {
            "state": 3
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-action-correlation-4",
      "connector_name": "test-connector-name-action-correlation-4",
      "source_type": "resource",
      "event_type": "comment",
      "component":  "test-component-action-correlation-4",
      "resource": "test-resource-action-correlation-4",
      "output": "test-output-action-correlation-4"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "component": "test-component-action-correlation-4",
                  "connector": "test-connector-action-correlation-4",
                  "connector_name": "test-connector-name-action-correlation-4",
                  "resource": "test-resource-action-correlation-4"
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
        }
      }
    ]
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 0)._id }}",
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
                "_t": "stateinc",
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
      },
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
                "_t": "metaalarmattach"
              },
              {
                "_t": "comment",
                "m": "test-output-action-correlation-4"
              },
              {
                "_t": "changestate",
                "val": 3
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
