Feature: create and update meta alarm
  I need to be able to create and update meta alarm

  @concurrent
  Scenario: given meta alarm rule and 2 simultaneous events should update meta alarm
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-second-1
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-second-1",
      "connector_name": "test-connector-name-axe-correlation-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-second-1",
      "resource": "test-resource-axe-correlation-second-1-1",
      "state": 1,
      "output": "test-output-axe-correlation-second-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "meta": "{{ .metaAlarmRuleID }}"
          }
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-second-1",
        "connector_name": "test-connector-name-axe-correlation-second-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-second-1",
        "resource": "test-resource-axe-correlation-second-1-2",
        "state": 2,
        "output": "test-output-axe-correlation-second-1"
      },
      {
        "connector": "test-connector-axe-correlation-second-1",
        "connector_name": "test-connector-name-axe-correlation-second-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-second-1",
        "resource": "test-resource-axe-correlation-second-1-3",
        "state": 3,
        "output": "test-output-axe-correlation-second-1"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-second-1; Count: 3; Children: test-component-axe-correlation-second-1",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "state": {
              "_t": "stateinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
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
                  "component": "test-component-axe-correlation-second-1",
                  "connector": "test-connector-axe-correlation-second-1",
                  "connector_name": "test-connector-name-axe-correlation-second-1",
                  "resource": "test-resource-axe-correlation-second-1-1"
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-second-1",
                  "connector": "test-connector-axe-correlation-second-1",
                  "connector_name": "test-connector-name-axe-correlation-second-1",
                  "resource": "test-resource-axe-correlation-second-1-2"
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-second-1",
                  "connector": "test-connector-axe-correlation-second-1",
                  "connector_name": "test-connector-name-axe-correlation-second-1",
                  "resource": "test-resource-axe-correlation-second-1-3"
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
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
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
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
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
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
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
  Scenario: given meta alarm rule and 2 simultaneous events for the same entity should create meta alarm
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-second-2
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-second-2",
        "connector_name": "test-connector-name-axe-correlation-second-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-second-2",
        "resource": "test-resource-axe-correlation-second-2",
        "state": 2,
        "output": "test-output-axe-correlation-second-2"
      },
      {
        "connector": "test-connector-axe-correlation-second-2",
        "connector_name": "test-connector-name-axe-correlation-second-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-second-2",
        "resource": "test-resource-axe-correlation-second-2",
        "state": 3,
        "output": "test-output-axe-correlation-second-2"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-2&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-second-2; Count: 1; Children: test-component-axe-correlation-second-2",
            "children": [
              "test-resource-axe-correlation-second-2/test-component-axe-correlation-second-2"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "state": {
              "_t": "stateinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
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
                  "component": "test-component-axe-correlation-second-2",
                  "connector": "test-connector-axe-correlation-second-2",
                  "connector_name": "test-connector-name-axe-correlation-second-2",
                  "resource": "test-resource-axe-correlation-second-2"
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
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
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
              "total_count": 4
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given 2 meta alarm rules and event should create 2 meta alarms
    Given I am admin
    When I save response metaAlarmRuleID1=test-metaalarmrule-axe-correlation-second-3-1
    When I save response metaAlarmRuleID2=test-metaalarmrule-axe-correlation-second-3-2
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-second-3",
      "connector_name": "test-connector-name-axe-correlation-second-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-second-3",
      "resource": "test-resource-axe-correlation-second-3",
      "state": 3,
      "output": "test-output-axe-correlation-second-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-3&correlation=true&sort_by=v.output&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-second-3-1; Count: 1; Children: test-component-axe-correlation-second-3",
            "children": [
              "test-resource-axe-correlation-second-3/test-component-axe-correlation-second-3"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID1 }}",
            "state": {
              "_t": "stateinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 1
            }
          }
        },
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-second-3-2; Count: 1; Children: test-component-axe-correlation-second-3",
            "children": [
              "test-resource-axe-correlation-second-3/test-component-axe-correlation-second-3"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID2 }}",
            "state": {
              "_t": "stateinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
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
    When I save response metaAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID1 }}",
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
                  "component": "test-component-axe-correlation-second-3",
                  "connector": "test-connector-axe-correlation-second-3",
                  "connector_name": "test-connector-name-axe-correlation-second-3",
                  "resource": "test-resource-axe-correlation-second-3"
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
        "_id": "{{ .metaAlarmID2 }}",
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
                  "component": "test-component-axe-correlation-second-3",
                  "connector": "test-connector-axe-correlation-second-3",
                  "connector_name": "test-connector-name-axe-correlation-second-3",
                  "resource": "test-resource-axe-correlation-second-3"
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
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
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
  Scenario: given 2 meta alarm rules and event should update 2 meta alarms
    Given I am admin
    When I save response metaAlarmRuleID1=test-metaalarmrule-axe-correlation-second-4-1
    When I save response metaAlarmRuleID2=test-metaalarmrule-axe-correlation-second-4-2
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-second-4",
      "connector_name": "test-connector-name-axe-correlation-second-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-second-4",
      "resource": "test-resource-axe-correlation-second-4-1",
      "state": 2,
      "output": "test-output-axe-correlation-second-4"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-4&correlation=true&sort_by=v.output&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "meta": "{{ .metaAlarmRuleID1 }}"
          }
        },
        {
          "v": {
            "meta": "{{ .metaAlarmRuleID2 }}"
          }
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-second-4",
      "connector_name": "test-connector-name-axe-correlation-second-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-second-4",
      "resource": "test-resource-axe-correlation-second-4-2",
      "state": 3,
      "output": "test-output-axe-correlation-second-4"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-4&correlation=true&sort_by=v.output&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-second-4-1; Count: 2; Children: test-component-axe-correlation-second-4",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID1 }}",
            "state": {
              "_t": "stateinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 1
            }
          }
        },
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-second-4-2; Count: 2; Children: test-component-axe-correlation-second-4",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID2 }}",
            "state": {
              "_t": "stateinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
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
    When I save response metaAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID1 }}",
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
                  "component": "test-component-axe-correlation-second-4",
                  "connector": "test-connector-axe-correlation-second-4",
                  "connector_name": "test-connector-name-axe-correlation-second-4",
                  "resource": "test-resource-axe-correlation-second-4-1"
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-second-4",
                  "connector": "test-connector-axe-correlation-second-4",
                  "connector_name": "test-connector-name-axe-correlation-second-4",
                  "resource": "test-resource-axe-correlation-second-4-2"
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
        }
      }
    ]
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID2 }}",
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
                  "component": "test-component-axe-correlation-second-4",
                  "connector": "test-connector-axe-correlation-second-4",
                  "connector_name": "test-connector-name-axe-correlation-second-4",
                  "resource": "test-resource-axe-correlation-second-4-1"
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-second-4",
                  "connector": "test-connector-axe-correlation-second-4",
                  "connector_name": "test-connector-name-axe-correlation-second-4",
                  "resource": "test-resource-axe-correlation-second-4-2"
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
        }
      }
    ]
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
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
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
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
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
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
  Scenario: given meta alarm event and child event should update child alarm in correct order
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-second-5
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-second-5",
      "connector_name": "test-connector-name-axe-correlation-second-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-second-5",
      "resource": "test-resource-axe-correlation-second-5",
      "state": 3,
      "output": "test-output-axe-correlation-second-5"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-5&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "meta": "{{ .metaAlarmRuleID }}"
          }
        }
      ]
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I send an event:
    """json
    [
      {
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "source_type": "resource",
        "event_type": "snooze",
        "duration": 6000,
        "component": "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "output": "test-output-axe-correlation-8",
        "initiator": "user"
      },
      {
        "connector": "test-connector-axe-correlation-second-5",
        "connector_name": "test-connector-name-axe-correlation-second-5",
        "source_type": "resource",
        "event_type": "ack",
        "component": "test-component-axe-correlation-second-5",
        "resource": "test-resource-axe-correlation-second-5",
        "output": "test-output-axe-correlation-second-5",
        "initiator": "user"
      }
    ]
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "snooze",
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "component": "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "source_type": "resource"
      },
      {
        "event_type": "snooze",
        "connector": "test-connector-axe-correlation-second-5",
        "connector_name": "test-connector-name-axe-correlation-second-5",
        "component": "test-component-axe-correlation-second-5",
        "resource": "test-resource-axe-correlation-second-5",
        "source_type": "resource"
      },
      {
        "event_type": "ack",
        "connector": "test-connector-axe-correlation-second-5",
        "connector_name": "test-connector-name-axe-correlation-second-5",
        "component": "test-component-axe-correlation-second-5",
        "resource": "test-resource-axe-correlation-second-5",
        "source_type": "resource"
      }
    ]
    """
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
                  "component": "test-component-axe-correlation-second-5",
                  "connector": "test-connector-axe-correlation-second-5",
                  "connector_name": "test-connector-name-axe-correlation-second-5",
                  "resource": "test-resource-axe-correlation-second-5"
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
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
              },
              {
                "_t": "ack",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user"
              },
              {
                "_t": "snooze",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user"
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
  Scenario: given manual meta alarm and updated or removed child alarm should update meta alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-second-6",
        "connector_name": "test-connector-name-axe-correlation-second-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-second-6",
        "resource": "test-resource-axe-correlation-second-6-1",
        "state": 1,
        "output": "test-output-axe-correlation-second-6"
      },
      {
        "connector": "test-connector-axe-correlation-second-6",
        "connector_name": "test-connector-name-axe-correlation-second-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-second-6",
        "resource": "test-resource-axe-correlation-second-6-2",
        "state": 2,
        "output": "test-output-axe-correlation-second-6"
      },
      {
        "connector": "test-connector-axe-correlation-second-6",
        "connector_name": "test-connector-name-axe-correlation-second-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-second-6",
        "resource": "test-resource-axe-correlation-second-6-3",
        "state": 3,
        "output": "test-output-axe-correlation-second-6"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-6&sort_by=v.resource&sort=asc
    When I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    When I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I save response alarmId3={{ (index .lastResponse.data 2)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaAlarm-axe-correlation-second-6",
      "comment": "test-metaAlarm-axe-correlation-second-6-comment",
      "alarms": ["{{ .alarmId1 }}", "{{ .alarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-6-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "display_name": "test-metaAlarm-axe-correlation-second-6"
          }
        }
      ]
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-6&correlation=true&sort_by=v.meta&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "output": "test-metaAlarm-axe-correlation-second-6-comment",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "display_name": "test-metaAlarm-axe-correlation-second-6",
            "state": {
              "_t": "stateinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 2
            },
            "status": {
              "_t": "statusinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 1
            }
          }
        },
        {
          "v": {
            "resource": "test-resource-axe-correlation-second-6-3"
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
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmEntityID={{ (index .lastResponse.data 0).entity._id }}
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
                  "component": "test-component-axe-correlation-second-6",
                  "connector": "test-connector-axe-correlation-second-6",
                  "connector_name": "test-connector-name-axe-correlation-second-6",
                  "resource": "test-resource-axe-correlation-second-6-1",
                  "parents": ["{{ .metaAlarmEntityID }}"]
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-second-6",
                  "connector": "test-connector-axe-correlation-second-6",
                  "connector_name": "test-connector-name-axe-correlation-second-6",
                  "resource": "test-resource-axe-correlation-second-6-2",
                  "parents": ["{{ .metaAlarmEntityID }}"]
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
        }
      }
    ]
    """
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ .metaAlarmID }}/add:
    """json
    {
      "comment": "test-metaAlarm-axe-correlation-second-6-comment",
      "alarms": ["{{ .alarmId3 }}"]
    }
    """
    Then the response code should be 204
    When I save request:
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
    When I do POST /api/v4/alarm-details until response code is 207 and body contains:
    """json
    [
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "component": "test-component-axe-correlation-second-6",
                  "connector": "test-connector-axe-correlation-second-6",
                  "connector_name": "test-connector-name-axe-correlation-second-6",
                  "resource": "test-resource-axe-correlation-second-6-1",
                  "parents": ["{{ .metaAlarmEntityID }}"]
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-second-6",
                  "connector": "test-connector-axe-correlation-second-6",
                  "connector_name": "test-connector-name-axe-correlation-second-6",
                  "resource": "test-resource-axe-correlation-second-6-2",
                  "parents": ["{{ .metaAlarmEntityID }}"]
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-second-6",
                  "connector": "test-connector-axe-correlation-second-6",
                  "connector_name": "test-connector-name-axe-correlation-second-6",
                  "resource": "test-resource-axe-correlation-second-6-3",
                  "parents": ["{{ .metaAlarmEntityID }}"]
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
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ .metaAlarmID }}/remove:
    """json
    {
      "comment": "test-metaAlarm-axe-correlation-second-6-comment",
      "alarms": ["{{ .alarmId3 }}"]
    }
    """
    Then the response code should be 204
    When I save request:
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
    When I do POST /api/v4/alarm-details until response code is 207 and body contains:
    """json
    [
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "component": "test-component-axe-correlation-second-6",
                  "connector": "test-connector-axe-correlation-second-6",
                  "connector_name": "test-connector-name-axe-correlation-second-6",
                  "resource": "test-resource-axe-correlation-second-6-1",
                  "parents": ["{{ .metaAlarmEntityID }}"]
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-second-6",
                  "connector": "test-connector-axe-correlation-second-6",
                  "connector_name": "test-connector-name-axe-correlation-second-6",
                  "resource": "test-resource-axe-correlation-second-6-2",
                  "parents": ["{{ .metaAlarmEntityID }}"]
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
        }
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-6-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-correlation-second-6",
            "connector": "test-connector-axe-correlation-second-6",
            "connector_name": "test-connector-name-axe-correlation-second-6",
            "resource": "test-resource-axe-correlation-second-6-3",
            "parents": []
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
  Scenario: given meta alarm and ack event should double ack children if they're already acked
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-axe-correlation-second-7",
      "type": "attribute",
      "output_template": "Count: {{ `{{ .Count }}` }}",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-axe-correlation-second-7"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-second-7",
        "connector_name": "test-connector-name-axe-correlation-second-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-second-7",
        "resource": "test-resource-axe-correlation-second-7-1",
        "state": 2,
        "output": "test-output-axe-correlation-second-7",
        "long_output": "test-long-output-axe-correlation-second-7",
        "author": "test-author-axe-correlation-second-7"
      },
      {
        "connector": "test-connector-axe-correlation-second-7",
        "connector_name": "test-connector-name-axe-correlation-second-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-second-7",
        "resource": "test-resource-axe-correlation-second-7-2",
        "state": 2,
        "output": "test-output-axe-correlation-second-7",
        "long_output": "test-long-output-axe-correlation-second-7",
        "author": "test-author-axe-correlation-second-7"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-7&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Count: 2",
            "meta": "{{ .metaAlarmRuleID }}"
          }
        }
      ]
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-second-7",
      "connector_name": "test-connector-name-axe-correlation-second-7",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-axe-correlation-second-7",
      "resource": "test-resource-axe-correlation-second-7-2",
      "output": "previous ack",
      "initiator": "user"
    }
    """
    When I do PUT /api/v4/alarms/{{ .metaAlarmID }}/ack:
    """json
    {
      "comment": "metaalarm ack"
    }
    """
    Then the response code should be 204
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ack",
        "connector": "api",
        "connector_name": "api",
        "component": "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "source_type": "resource"
      },
      {
        "event_type": "ack",
        "connector": "test-connector-axe-correlation-second-7",
        "connector_name": "test-connector-name-axe-correlation-second-7",
        "component": "test-component-axe-correlation-second-7",
        "resource": "test-resource-axe-correlation-second-7-1",
        "source_type": "resource"
      },
      {
        "event_type": "ack",
        "connector": "test-connector-axe-correlation-second-7",
        "connector_name": "test-connector-name-axe-correlation-second-7",
        "component": "test-component-axe-correlation-second-7",
        "resource": "test-resource-axe-correlation-second-7-2",
        "source_type": "resource"
      }
    ]
    """
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
                  "component": "test-component-axe-correlation-second-7",
                  "connector": "test-connector-axe-correlation-second-7",
                  "connector_name": "test-connector-name-axe-correlation-second-7",
                  "resource": "test-resource-axe-correlation-second-7-1"
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-second-7",
                  "connector": "test-connector-axe-correlation-second-7",
                  "connector_name": "test-connector-name-axe-correlation-second-7",
                  "resource": "test-resource-axe-correlation-second-7-2"
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
                "_t": "ack",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "metaalarm ack"
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
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
              },
              {
                "_t": "ack",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "metaalarm ack"
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
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
              },
              {
                "_t": "ack",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "previous ack"
              },
              {
                "_t": "ack",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "metaalarm ack"
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
  Scenario: given meta alarm child and cancel event should update parent state
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-axe-correlation-second-8-1",
      "type": "attribute",
      "auto_resolve": false,
      "output_template": "{{ `{{ .Rule.ID }}` }}",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-axe-correlation-second-8"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-axe-correlation-second-8-2",
      "type": "attribute",
      "auto_resolve": true,
      "output_template": "{{ `{{ .Rule.ID }}` }}",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-axe-correlation-second-8"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID2={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-second-8",
      "connector_name": "test-connector-name-axe-correlation-second-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-second-8",
      "resource": "test-resource-axe-correlation-second-8",
      "state": 2,
      "output": "test-output-axe-correlation-second-8"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-8&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "v": {
            "children": [
              "test-resource-axe-correlation-second-8/test-component-axe-correlation-second-8"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "state": {
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 2
            },
            "status": {
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 1
            }
          }
        },
        {
          "is_meta_alarm": true,
          "v": {
            "children": [
              "test-resource-axe-correlation-second-8/test-component-axe-correlation-second-8"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "state": {
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 2
            },
            "status": {
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-second-8",
      "connector_name": "test-connector-name-axe-correlation-second-8",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-axe-correlation-second-8",
      "resource": "test-resource-axe-correlation-second-8",
      "output": "test-output-axe-correlation-second-8"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-second-8",
      "connector_name": "test-connector-name-axe-correlation-second-8",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-axe-correlation-second-8",
      "resource": "test-resource-axe-correlation-second-8",
      "output": "test-output-axe-correlation-second-8"
    }
    """
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID1 }}&active_columns[]=v.output&correlation=true&opened=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "v": {
            "children": [
              "test-resource-axe-correlation-second-8/test-component-axe-correlation-second-8"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID1 }}",
            "state": {
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 0
            },
            "status": {
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
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
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID2 }}&active_columns[]=v.output&correlation=true&opened=false until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "v": {
            "children": [
              "test-resource-axe-correlation-second-8/test-component-axe-correlation-second-8"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID2 }}",
            "state": {
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 2
            },
            "status": {
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
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

  @concurrent
  Scenario: given meta alarm and assoc ticket event should ticket to children
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-axe-correlation-second-9",
      "type": "attribute",
      "auto_resolve": true,
      "output_template": "Count: {{ `{{ .Count }}` }}",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-axe-correlation-second-9"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-second-9",
        "connector_name": "test-connector-name-axe-correlation-second-9",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-second-9",
        "resource": "test-resource-axe-correlation-second-9-1",
        "state": 2,
        "output": "test-output-axe-correlation-second-9"
      },
      {
        "connector": "test-connector-axe-correlation-second-9",
        "connector_name": "test-connector-name-axe-correlation-second-9",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-second-9",
        "resource": "test-resource-axe-correlation-second-9-2",
        "state": 2,
        "output": "test-output-axe-correlation-second-9"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-9&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Count: 2",
            "meta": "{{ .metaAlarmRuleID }}"
          }
        }
      ]
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do PUT /api/v4/alarms/{{ .metaAlarmID }}/assocticket:
    """json
    {
      "ticket": "test-ticket-axe-correlation-second-9",
      "url": "test-url-axe-correlation-second-9",
      "system_name": "test-system-name-axe-correlation-second-9",
      "data": {
        "ticket_param_1": "ticket_value_1"
      },
      "comment": "test-comment-axe-correlation-second-9"
    }
    """
    Then the response code should be 204
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "assocticket",
        "connector": "api",
        "connector_name": "api",
        "component": "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "source_type": "resource"
      },
      {
        "event_type": "assocticket",
        "connector": "test-connector-axe-correlation-second-9",
        "connector_name": "test-connector-name-axe-correlation-second-9",
        "component": "test-component-axe-correlation-second-9",
        "resource": "test-resource-axe-correlation-second-9-1",
        "source_type": "resource"
      },
      {
        "event_type": "assocticket",
        "connector": "test-connector-axe-correlation-second-9",
        "connector_name": "test-connector-name-axe-correlation-second-9",
        "component": "test-component-axe-correlation-second-9",
        "resource": "test-resource-axe-correlation-second-9-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-9-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "assocticket",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "Ticket ID: test-ticket-axe-correlation-second-9. Ticket URL: test-url-axe-correlation-second-9. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-ticket-axe-correlation-second-9",
                "ticket_url": "test-url-axe-correlation-second-9",
                "ticket_system_name": "test-system-name-axe-correlation-second-9",
                "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-comment-axe-correlation-second-9"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "Ticket ID: test-ticket-axe-correlation-second-9. Ticket URL: test-url-axe-correlation-second-9. Ticket ticket_param_1: ticket_value_1.",
              "ticket": "test-ticket-axe-correlation-second-9",
              "ticket_url": "test-url-axe-correlation-second-9",
              "ticket_system_name": "test-system-name-axe-correlation-second-9",
              "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
              "ticket_data": {
                "ticket_param_1": "ticket_value_1"
              },
              "ticket_comment": "test-comment-axe-correlation-second-9"
            },
            "children": [],
            "component": "test-component-axe-correlation-second-9",
            "connector": "test-connector-axe-correlation-second-9",
            "connector_name": "test-connector-name-axe-correlation-second-9",
            "parents": [
              "{{ .metaAlarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-second-9-1",
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
                "_t": "metaalarmattach"
              },
              {
                "_t": "assocticket",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "Ticket ID: test-ticket-axe-correlation-second-9. Ticket URL: test-url-axe-correlation-second-9. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-ticket-axe-correlation-second-9",
                "ticket_url": "test-url-axe-correlation-second-9",
                "ticket_system_name": "test-system-name-axe-correlation-second-9",
                "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-comment-axe-correlation-second-9"
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
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-second-9-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "assocticket",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "Ticket ID: test-ticket-axe-correlation-second-9. Ticket URL: test-url-axe-correlation-second-9. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-ticket-axe-correlation-second-9",
                "ticket_url": "test-url-axe-correlation-second-9",
                "ticket_system_name": "test-system-name-axe-correlation-second-9",
                "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-comment-axe-correlation-second-9"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "Ticket ID: test-ticket-axe-correlation-second-9. Ticket URL: test-url-axe-correlation-second-9. Ticket ticket_param_1: ticket_value_1.",
              "ticket": "test-ticket-axe-correlation-second-9",
              "ticket_url": "test-url-axe-correlation-second-9",
              "ticket_system_name": "test-system-name-axe-correlation-second-9",
              "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
              "ticket_data": {
                "ticket_param_1": "ticket_value_1"
              },
              "ticket_comment": "test-comment-axe-correlation-second-9"
            },
            "children": [],
            "component": "test-component-axe-correlation-second-9",
            "connector": "test-connector-axe-correlation-second-9",
            "connector_name": "test-connector-name-axe-correlation-second-9",
            "parents": [
              "{{ .metaAlarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-second-9-2",
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
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
              },
              {
                "_t": "assocticket",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "Ticket ID: test-ticket-axe-correlation-second-9. Ticket URL: test-url-axe-correlation-second-9. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-ticket-axe-correlation-second-9",
                "ticket_url": "test-url-axe-correlation-second-9",
                "ticket_system_name": "test-system-name-axe-correlation-second-9",
                "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-comment-axe-correlation-second-9"
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
