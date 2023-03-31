Feature: Manual meta alarms

  @concurrent
  Scenario: given alarms should create manual meta alarms
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-manual-correlation-1",
        "connector_name": "test-connector-name-manual-correlation-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-manual-correlation-1",
        "resource": "test-resource-manual-correlation-1-1",
        "state": 1,
        "output": "test-output-manual-correlation-1"
      },
      {
        "connector": "test-connector-manual-correlation-1",
        "connector_name": "test-connector-name-manual-correlation-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-manual-correlation-1",
        "resource": "test-resource-manual-correlation-1-2",
        "state": 1,
        "output": "test-output-manual-correlation-1"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-manual-correlation-1&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-manual-correlation-1-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-manual-correlation-1-2"
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
    When I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    When I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-manual-correlation-1-1",
      "comment": "test-metaalarm-manual-correlation-1-1-comment",
      "alarms": ["{{ .alarmId1 }}"]
    }
    """
    Then the response code should be 204
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-manual-correlation-1-2",
      "comment": "test-metaalarm-manual-correlation-1-2-comment",
      "alarms": ["{{ .alarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/manual-meta-alarms?search=test-metaalarm-manual-correlation-1 until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-metaalarm-manual-correlation-1-1"
      },
      {
        "name": "test-metaalarm-manual-correlation-1-2"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-manual-correlation-1&correlation=true&sort_by=v.display_name&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "_id": "test-metaalarm-rule-manual",
            "name": "test-metaalarm-rule-manual-name",
            "type": "manualgroup"
          },
          "v": {
            "connector": "engine",
            "connector_name": "correlation",
            "component": "metaalarm",
            "display_name": "test-metaalarm-manual-correlation-1-1",
            "meta": "test-metaalarm-rule-manual",
            "children": [
              "test-resource-manual-correlation-1-1/test-component-manual-correlation-1"
            ],
            "output": "test-metaalarm-manual-correlation-1-1-comment"
          }
        },
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "_id": "test-metaalarm-rule-manual",
            "name": "test-metaalarm-rule-manual-name",
            "type": "manualgroup"
          },
          "v": {
            "connector": "engine",
            "connector_name": "correlation",
            "component": "metaalarm",
            "display_name": "test-metaalarm-manual-correlation-1-2",
            "meta": "test-metaalarm-rule-manual",
            "children": [
              "test-resource-manual-correlation-1-2/test-component-manual-correlation-1"
            ],
            "output": "test-metaalarm-manual-correlation-1-2-comment",
            "state": {
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
    When I save response metaAlarmId1={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmEntityId1={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmId2={{ (index .lastResponse.data 1)._id }}
    When I save response metaAlarmEntityId2={{ (index .lastResponse.data 1).entity._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmId1 }}",
        "children": {
          "page": 1
        },
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .metaAlarmId2 }}",
        "children": {
          "page": 1
        },
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .alarmId1 }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .alarmId2 }}",
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-connector-manual-correlation-1",
                  "connector_name": "test-connector-name-manual-correlation-1",
                  "component": "test-component-manual-correlation-1",
                  "resource": "test-resource-manual-correlation-1-1"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 1
            }
          },
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "m": "test-metaalarm-manual-correlation-1-1-comment",
                "val": 1
              },
              {
                "_t": "statusinc",
                "m": "test-metaalarm-manual-correlation-1-1-comment",
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
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-connector-manual-correlation-1",
                  "connector_name": "test-connector-name-manual-correlation-1",
                  "component": "test-component-manual-correlation-1",
                  "resource": "test-resource-manual-correlation-1-2"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 1
            }
          },
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "m": "test-metaalarm-manual-correlation-1-2-comment",
                "val": 1
              },
              {
                "_t": "statusinc",
                "m": "test-metaalarm-manual-correlation-1-2-comment",
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
                "_t": "statusinc"
              },
              {
                "_t": "metaalarmattach",
                "m": "Rule: {test-metaalarm-rule-manual-name}\n Displayname: {test-metaalarm-manual-correlation-1-1}\n Entity: {{ `{` }}{{ .metaAlarmEntityId1 }}{{ `}` }}"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "metaalarmattach",
                "m": "Rule: {test-metaalarm-rule-manual-name}\n Displayname: {test-metaalarm-manual-correlation-1-2}\n Entity: {{ `{` }}{{ .metaAlarmEntityId2 }}{{ `}` }}"
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
  Scenario: given manual meta alarm should add alarm to it
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-manual-correlation-2",
        "connector_name": "test-connector-name-manual-correlation-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-manual-correlation-2",
        "resource": "test-resource-manual-correlation-2-1",
        "state": 1,
        "output": "test-output-manual-correlation-2"
      },
      {
        "connector": "test-connector-manual-correlation-2",
        "connector_name": "test-connector-name-manual-correlation-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-manual-correlation-2",
        "resource": "test-resource-manual-correlation-2-2",
        "state": 1,
        "output": "test-output-manual-correlation-2"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-manual-correlation-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-manual-correlation-2-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-manual-correlation-2-2"
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
    When I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    When I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-manual-correlation-2",
      "comment": "test-metaalarm-manual-correlation-2-1-comment",
      "alarms": ["{{ .alarmId1 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/manual-meta-alarms?search=test-metaalarm-manual-correlation-2 until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-metaalarm-manual-correlation-2"
      }
    ]
    """
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ (index .lastResponse 0)._id }}/add:
    """json
    {
      "comment": "test-metaalarm-manual-correlation-2-2-comment",
      "alarms": ["{{ .alarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-manual-correlation-2&correlation=true&sort_by=v.display_name&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "_id": "test-metaalarm-rule-manual",
            "name": "test-metaalarm-rule-manual-name",
            "type": "manualgroup"
          },
          "v": {
            "display_name": "test-metaalarm-manual-correlation-2",
            "output": "test-metaalarm-manual-correlation-2-2-comment"
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
    When I save response metaAlarmId={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmEntityId={{ (index .lastResponse.data 0).entity._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmId }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        },
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .alarmId2 }}",
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-connector-manual-correlation-2",
                  "connector_name": "test-connector-name-manual-correlation-2",
                  "component": "test-component-manual-correlation-2",
                  "resource": "test-resource-manual-correlation-2-1"
                }
              },
              {
                "v": {
                  "connector": "test-connector-manual-correlation-2",
                  "connector_name": "test-connector-name-manual-correlation-2",
                  "component": "test-component-manual-correlation-2",
                  "resource": "test-resource-manual-correlation-2-2"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          },
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "m": "test-metaalarm-manual-correlation-2-1-comment",
                "val": 1
              },
              {
                "_t": "statusinc",
                "m": "test-metaalarm-manual-correlation-2-1-comment",
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
                "_t": "statusinc"
              },
              {
                "_t": "metaalarmattach",
                "m": "Rule: {test-metaalarm-rule-manual-name}\n Displayname: {test-metaalarm-manual-correlation-2}\n Entity: {{ `{` }}{{ .metaAlarmEntityId }}{{ `}` }}"
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
  Scenario: given manual meta alarm should remove alarm from it
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-manual-correlation-3",
        "connector_name": "test-connector-name-manual-correlation-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-manual-correlation-3",
        "resource": "test-resource-manual-correlation-3-1",
        "state": 1,
        "output": "test-output-manual-correlation-3"
      },
      {
        "connector": "test-connector-manual-correlation-3",
        "connector_name": "test-connector-name-manual-correlation-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-manual-correlation-3",
        "resource": "test-resource-manual-correlation-3-2",
        "state": 1,
        "output": "test-output-manual-correlation-3"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-manual-correlation-3&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-manual-correlation-3-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-manual-correlation-3-2"
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
    When I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    When I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-manual-correlation-3",
      "comment": "test-metaalarm-manual-correlation-3-1-comment",
      "alarms": [
        "{{ .alarmId1 }}",
        "{{ .alarmId2 }}"
      ]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/manual-meta-alarms?search=test-metaalarm-manual-correlation-3 until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-metaalarm-manual-correlation-3"
      }
    ]
    """
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ (index .lastResponse 0)._id }}/remove:
    """json
    {
      "comment": "test-metaalarm-manual-correlation-3-2-comment",
      "alarms": ["{{ .alarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-manual-correlation-3&correlation=true&sort_by=v.component&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "_id": "test-metaalarm-rule-manual",
            "name": "test-metaalarm-rule-manual-name",
            "type": "manualgroup"
          },
          "v": {
            "display_name": "test-metaalarm-manual-correlation-3",
            "output": "test-metaalarm-manual-correlation-3-2-comment"
          }
        },
        {
          "v": {
            "resource": "test-resource-manual-correlation-3-2"
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
    When I save response metaAlarmId={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmEntityId={{ (index .lastResponse.data 0).entity._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmId }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        },
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .alarmId2 }}",
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-connector-manual-correlation-3",
                  "connector_name": "test-connector-name-manual-correlation-3",
                  "component": "test-component-manual-correlation-3",
                  "resource": "test-resource-manual-correlation-3-1"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 1
            }
          },
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "m": "test-metaalarm-manual-correlation-3-1-comment",
                "val": 1
              },
              {
                "_t": "statusinc",
                "m": "test-metaalarm-manual-correlation-3-1-comment",
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
                "_t": "statusinc"
              },
              {
                "_t": "metaalarmattach",
                "m": "Rule: {test-metaalarm-rule-manual-name}\n Displayname: {test-metaalarm-manual-correlation-3}\n Entity: {{ `{` }}{{ .metaAlarmEntityId }}{{ `}` }}"
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
  Scenario: given get manual meta alarm request with not exist search should return empty response
    When I am admin
    When I do GET /api/v4/cat/manual-meta-alarms?search=not-exist
    Then the response code should be 200
    Then the response body should be:
    """json
    []
    """
