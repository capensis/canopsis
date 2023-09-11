Feature: Get alarms
  I need to be able to get a alarms

  @concurrent
  Scenario: given get correlation request should return meta alarms or alarms without parent
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-correlation-alarm-get-1",
        "connector_name": "test-connector-name-correlation-alarm-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-1",
        "resource": "test-resource-correlation-alarm-get-1-1",
        "state": 1,
        "output": "test-output-correlation-alarm-get-1"
      },
      {
        "connector": "test-connector-correlation-alarm-get-1",
        "connector_name": "test-connector-name-correlation-alarm-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-1",
        "resource": "test-resource-correlation-alarm-get-1-2",
        "state": 1,
        "output": "test-output-correlation-alarm-get-1"
      },
      {
        "connector": "test-connector-correlation-alarm-get-1",
        "connector_name": "test-connector-name-correlation-alarm-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-1",
        "resource": "test-resource-correlation-alarm-get-1-3",
        "state": 1,
        "output": "test-output-correlation-alarm-get-1"
      },
      {
        "connector": "test-connector-correlation-alarm-get-1",
        "connector_name": "test-connector-name-correlation-alarm-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-1",
        "resource": "test-resource-correlation-alarm-get-1-4",
        "state": 1,
        "output": "test-output-correlation-alarm-get-1"
      }
    ]
    """
    When I do GET /api/v4/alarms?correlation=false&search=test-resource-correlation-alarm-get-1&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-correlation-alarm-get-1-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-correlation-alarm-get-1-2"
          }
        },
        {
          "v": {
            "resource": "test-resource-correlation-alarm-get-1-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-correlation-alarm-get-1-4"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """
    When I save response childAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response childAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response childAlarmID3={{ (index .lastResponse.data 2)._id }}
    When I do GET /api/v4/alarms?correlation=true&search=test-resource-correlation-alarm-get-1&multi_sort[]=v.component,asc&multi_sort[]=meta_alarm_rule.name,asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "_id": "test-metaalarm-rule-alarm-correlation-get-1-1",
            "name": "test-metaalarm-rule-alarm-correlation-get-1-1-name"
          },
          "v": {
            "meta": "test-metaalarm-rule-alarm-correlation-get-1-1",
            "output": "Rule: test-metaalarm-rule-alarm-correlation-get-1-1; Count: 3; Children: test-component-correlation-alarm-get-1"
          }
        },
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "_id": "test-metaalarm-rule-alarm-correlation-get-1-2",
            "name": "test-metaalarm-rule-alarm-correlation-get-1-2-name"
          },
          "v": {
            "meta": "test-metaalarm-rule-alarm-correlation-get-1-2",
            "output": "Rule: test-metaalarm-rule-alarm-correlation-get-1-2; Count: 1; Children: test-component-correlation-alarm-get-1",
            "children": [
              "test-resource-correlation-alarm-get-1-3/test-component-correlation-alarm-get-1"
            ]
          }
        },
        {
          "is_meta_alarm": false,
          "v": {
            "resource": "test-resource-correlation-alarm-get-1-4"
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
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
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
      },
      {
        "_id": "{{ .metaAlarmID }}",
        "children": {
          "page": 2,
          "limit": 2,
          "multi_sort": ["v.resource,asc"]
        }
      },
      {
        "_id": "{{ .metaAlarmID }}",
        "steps": {
          "page": 1
        },
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ .childAlarmID1 }}",
        "steps": {
          "page": 1
        },
        "children": {
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
                  "resource": "test-resource-correlation-alarm-get-1-1"
                },
                "parents": 1,
                "meta_alarm_rules": [
                  {
                    "_id": "test-metaalarm-rule-alarm-correlation-get-1-1",
                    "name": "test-metaalarm-rule-alarm-correlation-get-1-1-name"
                  }
                ]
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-1-2"
                },
                "parents": 1,
                "meta_alarm_rules": [
                  {
                    "_id": "test-metaalarm-rule-alarm-correlation-get-1-1",
                    "name": "test-metaalarm-rule-alarm-correlation-get-1-1-name"
                  }
                ]
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-1-3"
                },
                "parents": 2,
                "meta_alarm_rules": [
                  {
                    "_id": "test-metaalarm-rule-alarm-correlation-get-1-1",
                    "name": "test-metaalarm-rule-alarm-correlation-get-1-1-name"
                  },
                  {
                    "_id": "test-metaalarm-rule-alarm-correlation-get-1-2",
                    "name": "test-metaalarm-rule-alarm-correlation-get-1-2-name"
                  }
                ]
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
          "children": {
            "data": [
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-1-3"
                }
              }
            ],
            "meta": {
              "page": 2,
              "page_count": 2,
              "per_page": 2,
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
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          },
          "children": {
            "data": [
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-1-1"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-1-2"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-1-3"
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
                "_t": "metaalarmattach"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          },
          "children": {
            "data": [],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 0
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given get correlation request should return resolved children
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-correlation-alarm-get-2",
        "connector_name": "test-connector-name-correlation-alarm-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-2",
        "resource": "test-resource-correlation-alarm-get-2-1",
        "state": 1,
        "output": "test-output-correlation-alarm-get-2"
      },
      {
        "connector": "test-connector-correlation-alarm-get-2",
        "connector_name": "test-connector-name-correlation-alarm-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-2",
        "resource": "test-resource-correlation-alarm-get-2-2",
        "state": 1,
        "output": "test-output-correlation-alarm-get-2"
      },
      {
        "connector": "test-connector-correlation-alarm-get-2",
        "connector_name": "test-connector-name-correlation-alarm-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-2",
        "resource": "test-resource-correlation-alarm-get-2-3",
        "state": 1,
        "output": "test-output-correlation-alarm-get-2"
      }
    ]
    """
    When I do GET /api/v4/alarms?correlation=false&search=test-resource-correlation-alarm-get-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response childAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response childAlarmID3={{ (index .lastResponse.data 2)._id }}
    When I do GET /api/v4/alarms?correlation=true&search=test-resource-correlation-alarm-get-2&multi_sort[]=v.component,asc&multi_sort[]=meta_alarm_rule.name,asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 3
        },
        {
          "children": 2
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
    When I save response metaAlarmResource2={{ (index .lastResponse.data 1).v.resource }}
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-correlation-alarm-get-2",
        "connector_name": "test-connector-name-correlation-alarm-get-2",
        "source_type": "resource",
        "event_type": "cancel",
        "component": "test-component-correlation-alarm-get-2",
        "resource": "test-resource-correlation-alarm-get-2-2"
      },
      {
        "connector": "test-connector-correlation-alarm-get-2",
        "connector_name": "test-connector-name-correlation-alarm-get-2",
        "source_type": "resource",
        "event_type": "cancel",
        "component": "test-component-correlation-alarm-get-2",
        "resource": "test-resource-correlation-alarm-get-2-3"
      },
      {
        "connector": "engine",
        "connector_name": "correlation",
        "source_type": "resource",
        "event_type": "cancel",
        "component": "metaalarm",
        "resource": "{{ .metaAlarmResource2 }}"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-correlation-alarm-get-2",
        "connector_name": "test-connector-name-correlation-alarm-get-2",
        "source_type": "resource",
        "event_type": "resolve_cancel",
        "component": "test-component-correlation-alarm-get-2",
        "resource": "test-resource-correlation-alarm-get-2-2"
      },
      {
        "connector": "test-connector-correlation-alarm-get-2",
        "connector_name": "test-connector-name-correlation-alarm-get-2",
        "source_type": "resource",
        "event_type": "resolve_cancel",
        "component": "test-component-correlation-alarm-get-2",
        "resource": "test-resource-correlation-alarm-get-2-3"
      },
      {
        "connector": "engine",
        "connector_name": "correlation",
        "source_type": "resource",
        "event_type": "resolve_cancel",
        "component": "metaalarm",
        "resource": "{{ .metaAlarmResource2 }}"
      }
    ]
    """
    When I do GET /api/v4/alarms?correlation=true&opened=true&search=test-resource-correlation-alarm-get-2&multi_sort[]=v.component,asc&multi_sort[]=meta_alarm_rule.name,asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 3
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
    When I do GET /api/v4/alarms?correlation=true&opened=false&search=test-resource-correlation-alarm-get-2&multi_sort[]=v.component,asc&multi_sort[]=meta_alarm_rule.name,asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 2
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
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
      },
      {
        "_id": "{{ .metaAlarmID1 }}",
        "opened": true,
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ .metaAlarmID2 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ .metaAlarmID2 }}",
        "opened": false,
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
                  "resource": "test-resource-correlation-alarm-get-2-1"
                },
                "parents": 2,
                "meta_alarm_rules": [
                  {
                    "_id": "test-metaalarm-rule-alarm-correlation-get-2-1",
                    "name": "test-metaalarm-rule-alarm-correlation-get-2-1-name"
                  },
                  {
                    "_id": "test-metaalarm-rule-alarm-correlation-get-2-2",
                    "name": "test-metaalarm-rule-alarm-correlation-get-2-2-name"
                  }
                ]
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-2-2"
                },
                "parents": 1,
                "meta_alarm_rules": [
                  {
                    "_id": "test-metaalarm-rule-alarm-correlation-get-2-1",
                    "name": "test-metaalarm-rule-alarm-correlation-get-2-1-name"
                  }
                ]
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-2-3"
                },
                "parents": 2,
                "meta_alarm_rules": [
                  {
                    "_id": "test-metaalarm-rule-alarm-correlation-get-2-1",
                    "name": "test-metaalarm-rule-alarm-correlation-get-2-1-name"
                  },
                  {
                    "_id": "test-metaalarm-rule-alarm-correlation-get-2-2",
                    "name": "test-metaalarm-rule-alarm-correlation-get-2-2-name"
                  }
                ]
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
          "children": {
            "data": [
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-2-1"
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
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-2-1"
                },
                "parents": 2,
                "meta_alarm_rules": [
                  {
                    "_id": "test-metaalarm-rule-alarm-correlation-get-2-1",
                    "name": "test-metaalarm-rule-alarm-correlation-get-2-1-name"
                  },
                  {
                    "_id": "test-metaalarm-rule-alarm-correlation-get-2-2",
                    "name": "test-metaalarm-rule-alarm-correlation-get-2-2-name"
                  }
                ]
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-2-3"
                },
                "parents": 2,
                "meta_alarm_rules": [
                  {
                    "_id": "test-metaalarm-rule-alarm-correlation-get-2-1",
                    "name": "test-metaalarm-rule-alarm-correlation-get-2-1-name"
                  },
                  {
                    "_id": "test-metaalarm-rule-alarm-correlation-get-2-2",
                    "name": "test-metaalarm-rule-alarm-correlation-get-2-2-name"
                  }
                ]
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
                  "resource": "test-resource-correlation-alarm-get-2-3"
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

  @concurrent
  Scenario: given get search correlation request should return filtered children
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-correlation-alarm-get-3",
        "connector_name": "test-connector-name-correlation-alarm-get-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-3",
        "resource": "test-resource-correlation-alarm-get-3-1",
        "state": 1,
        "output": "test-output-correlation-alarm-get-3"
      },
      {
        "connector": "test-connector-correlation-alarm-get-3",
        "connector_name": "test-connector-name-correlation-alarm-get-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-3",
        "resource": "test-resource-correlation-alarm-get-3-2",
        "state": 1,
        "output": "test-resource-correlation-alarm-get-3-search"
      },
      {
        "connector": "test-connector-correlation-alarm-get-3",
        "connector_name": "test-connector-name-correlation-alarm-get-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-3",
        "resource": "test-resource-correlation-alarm-get-3-3",
        "state": 1,
        "output": "test-resource-correlation-alarm-get-3-search"
      },
      {
        "connector": "test-connector-correlation-alarm-get-3",
        "connector_name": "test-connector-name-correlation-alarm-get-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-3",
        "resource": "test-resource-correlation-alarm-get-3-4",
        "state": 1,
        "output": "test-resource-correlation-alarm-get-3-search"
      },
      {
        "connector": "test-connector-correlation-alarm-get-3",
        "connector_name": "test-connector-name-correlation-alarm-get-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-3",
        "resource": "test-resource-correlation-alarm-get-3-5",
        "state": 1,
        "output": "test-resource-correlation-alarm-get-3-search"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-3&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response childAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response childAlarmID3={{ (index .lastResponse.data 2)._id }}
    When I save response childAlarmID4={{ (index .lastResponse.data 3)._id }}
    When I save response childAlarmID5={{ (index .lastResponse.data 4)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metalarm-correlation-alarm-get-3-1",
      "comment": "test-metalarm-correlation-alarm-get-3-1-comment",
      "alarms": [
        "{{ .childAlarmID1 }}",
        "{{ .childAlarmID2 }}",
        "{{ .childAlarmID3 }}"
      ]
    }
    """
    Then the response code should be 204
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metalarm-correlation-alarm-get-3-2",
      "comment": "test-metalarm-correlation-alarm-get-3-2-comment",
      "alarms": [
        "{{ .childAlarmID4 }}"
      ]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-3&correlation=true&sort_by=v.output&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 3
        },
        {
          "children": 1
        },
        {}
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-correlation-alarm-get-3",
      "connector_name": "test-connector-name-correlation-alarm-get-3",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-correlation-alarm-get-3",
      "resource": "test-resource-correlation-alarm-get-3-3"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-correlation-alarm-get-3",
      "connector_name": "test-connector-name-correlation-alarm-get-3",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-correlation-alarm-get-3",
      "resource": "test-resource-correlation-alarm-get-3-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-3-search&active_columns[]=v.output&correlation=true&sort_by=children&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3
        },
        {
          "is_meta_alarm": true,
          "children": 1
        },
        {
          "is_meta_alarm": false,
          "v": {
            "resource": "test-resource-correlation-alarm-get-3-5"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "search": "test-resource-correlation-alarm-get-3-search",
        "search_by": ["v.output"],
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ (index .lastResponse.data 1)._id }}",
        "search": "test-resource-correlation-alarm-get-3-search",
        "search_by": ["v.output"],
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
                  "resource": "test-resource-correlation-alarm-get-3-1"
                },
                "filtered": false
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-3-2"
                },
                "filtered": true
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-3-3"
                },
                "filtered": true
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
          "children": {
            "data": [
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-3-4"
                },
                "filtered": true
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
    When I do GET /api/v4/alarms?opened=true&search=test-resource-correlation-alarm-get-3-search&active_columns[]=v.output&correlation=true&sort_by=children&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3
        },
        {
          "is_meta_alarm": true,
          "children": 1
        },
        {
          "is_meta_alarm": false,
          "v": {
            "resource": "test-resource-correlation-alarm-get-3-5"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "opened": true,
        "search": "test-resource-correlation-alarm-get-3-search",
        "search_by": ["v.output"],
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ (index .lastResponse.data 1)._id }}",
        "opened": true,
        "search": "test-resource-correlation-alarm-get-3-search",
        "search_by": ["v.output"],
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
                  "resource": "test-resource-correlation-alarm-get-3-1"
                },
                "filtered": false
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-3-2"
                },
                "filtered": true
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
                  "resource": "test-resource-correlation-alarm-get-3-4"
                },
                "filtered": true
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
    When I do GET /api/v4/alarms?opened=true&search=test-resource-correlation-alarm-get-3-3&active_columns[]=v.resource&correlation=true&sort_by=children&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/alarms?opened=false&search=test-resource-correlation-alarm-get-3-3&active_columns[]=v.resource&correlation=true&sort_by=children&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "total_count": 0
      }
    }
    """

  @concurrent
  Scenario: given get correlation with_instructions request should return
  meta alarms with children, children should have assigned instruction if they have it, the corresponding
  metaalarm should have a mark about it
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-correlation-alarm-get-4",
        "connector_name": "test-connector-name-correlation-alarm-get-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-4",
        "resource": "test-resource-correlation-alarm-get-4-1",
        "state": 1,
        "output": "test-output-correlation-alarm-get-4"
      },
      {
        "connector": "test-connector-correlation-alarm-get-4",
        "connector_name": "test-connector-name-correlation-alarm-get-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-4",
        "resource": "test-resource-correlation-alarm-get-4-2",
        "state": 1,
        "output": "test-output-correlation-alarm-get-4"
      },
      {
        "connector": "test-connector-correlation-alarm-get-4",
        "connector_name": "test-connector-name-correlation-alarm-get-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-4",
        "resource": "test-resource-correlation-alarm-get-4-3",
        "state": 1,
        "output": "test-output-correlation-alarm-get-4"
      },
      {
        "connector": "test-connector-correlation-alarm-get-4",
        "connector_name": "test-connector-name-correlation-alarm-get-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-4",
        "resource": "test-resource-correlation-alarm-get-4-4",
        "state": 1,
        "output": "test-output-correlation-alarm-get-4"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-4&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response childAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response childAlarmID3={{ (index .lastResponse.data 2)._id }}
    When I save response childAlarmID4={{ (index .lastResponse.data 3)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metalarm-correlation-alarm-get-4-1",
      "comment": "test-metalarm-correlation-alarm-get-4-1-comment",
      "alarms": [
        "{{ .childAlarmID1 }}",
        "{{ .childAlarmID2 }}"
      ]
    }
    """
    Then the response code should be 204
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metalarm-correlation-alarm-get-4-2",
      "comment": "test-metalarm-correlation-alarm-get-4-2-comment",
      "alarms": [
        "{{ .childAlarmID3 }}",
        "{{ .childAlarmID4 }}"
      ]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-4&correlation=true&with_instructions=true&sort_by=v.display_name&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "children_instructions": true,
          "v": {
            "display_name": "test-metalarm-correlation-alarm-get-4-1"
          }
        },
        {
          "is_meta_alarm": true,
          "children": 2,
          "children_instructions": false,
          "v": {
            "display_name": "test-metalarm-correlation-alarm-get-4-2"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "with_instructions": true,
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ (index .lastResponse.data 1)._id }}",
        "with_instructions": true,
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
                  "resource": "test-resource-correlation-alarm-get-4-1",
                  "state": {
                    "val": 1
                  }
                },
                "assigned_instructions": []
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-4-2"
                },
                "assigned_instructions": [
                  {
                    "_id": "test-instruction-correlation-alarm-get-instruction-get-4-1"
                  }
                ]
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
                  "resource": "test-resource-correlation-alarm-get-4-3"
                },
                "assigned_instructions": []
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-4-4"
                },
                "assigned_instructions": []
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
  Scenario: given get correlation alarms with the same children, but different alarms shouldn't have alarms of each other
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-correlation-alarm-get-5",
        "connector_name": "test-connector-name-correlation-alarm-get-5",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-5",
        "resource": "test-resource-correlation-alarm-get-5-1",
        "state": 1,
        "output": "test-output-correlation-alarm-get-5"
      },
      {
        "connector": "test-connector-correlation-alarm-get-5",
        "connector_name": "test-connector-name-correlation-alarm-get-5",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-5",
        "resource": "test-resource-correlation-alarm-get-5-2",
        "state": 1,
        "output": "test-output-correlation-alarm-get-5"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-5&opened=true&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response childAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metalarm-correlation-alarm-get-5-1",
      "comment": "test-metalarm-correlation-alarm-get-5-1-comment",
      "alarms": [
        "{{ .childAlarmID1 }}",
        "{{ .childAlarmID2 }}"
      ]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-5&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 2
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-correlation-alarm-get-5",
      "connector_name": "test-connector-name-correlation-alarm-get-5",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-correlation-alarm-get-5",
      "resource": "test-resource-correlation-alarm-get-5-1",
      "output": "test-output-correlation-alarm-get-5"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-correlation-alarm-get-5",
      "connector_name": "test-connector-name-correlation-alarm-get-5",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-correlation-alarm-get-5",
      "resource": "test-resource-correlation-alarm-get-5-1",
      "output": "test-output-correlation-alarm-get-5"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-correlation-alarm-get-5",
      "connector_name": "test-connector-name-correlation-alarm-get-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-correlation-alarm-get-5",
      "resource": "test-resource-correlation-alarm-get-5-1",
      "state": 2,
      "output": "test-output-correlation-alarm-get-5"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-5&opened=true&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID3={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metalarm-correlation-alarm-get-5-2",
      "comment": "test-metalarm-correlation-alarm-get-5-2-comment",
      "alarms": [
        "{{ .childAlarmID3 }}",
        "{{ .childAlarmID2 }}"
      ]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-5&correlation=true&sort_by=v.display_name&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 2
        },
        {
          "children": 2
        }
      ]
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
      },
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
                "_id": "{{ .childAlarmID1 }}",
                "v": {
                  "resource": "test-resource-correlation-alarm-get-5-1",
                  "state": {
                    "val": 1
                  }
                }
              },
              {
                "_id": "{{ .childAlarmID2 }}",
                "v": {
                  "resource": "test-resource-correlation-alarm-get-5-2"
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
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "_id": "{{ .childAlarmID3 }}",
                "v": {
                  "resource": "test-resource-correlation-alarm-get-5-1",
                  "state": {
                    "val": 2
                  }
                }
              },
              {
                "_id": "{{ .childAlarmID2 }}",
                "v": {
                  "resource": "test-resource-correlation-alarm-get-5-2"
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

  @concurrent
  Scenario: given alarms which are matched to metaalarm should return instruction statuses and icons
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-1",
        "state": 1,
        "output": "test-output-correlation-alarm-get-6"
      },
      {
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-2",
        "state": 1,
        "output": "test-output-correlation-alarm-get-6"
      },
      {
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-3",
        "state": 1,
        "output": "test-output-correlation-alarm-get-6"
      },
      {
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-4",
        "state": 1,
        "output": "test-output-correlation-alarm-get-6"
      },
      {
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-5",
        "state": 1,
        "output": "test-output-correlation-alarm-get-6"
      }
    ]
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-2",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-3",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-4",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-5",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-3",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-3",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?correlation=false&search=test-resource-correlation-alarm-get-6&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response childAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response childAlarmID3={{ (index .lastResponse.data 2)._id }}
    When I save response childAlarmID4={{ (index .lastResponse.data 3)._id }}
    When I save response childAlarmID5={{ (index .lastResponse.data 4)._id }}
    When I do GET /api/v4/alarms?correlation=true&search=test-resource-correlation-alarm-get-6&multi_sort[]=v.component,asc&multi_sort[]=meta_alarm_rule.name,asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 5
        }
      ]
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .childAlarmID4 }}",
      "instruction": "test-instruction-correlation-alarm-get-6-4"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .childAlarmID5 }}",
      "instruction": "test-instruction-correlation-alarm-get-6-5"
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-4",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-4",
        "source_type": "resource"
      },
      {
        "event_type": "instructionstarted",
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-5",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-correlation-alarm-get-6",
        "connector_name": "test-connector-name-correlation-alarm-get-6",
        "component": "test-component-correlation-alarm-get-6",
        "resource": "test-resource-correlation-alarm-get-6-5",
        "source_type": "resource"
      }
    ]
    """
    When I save request:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "with_instructions": true,
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
                  "resource": "test-resource-correlation-alarm-get-6-1"
                },
                "instruction_execution_icon": 9
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-6-2"
                },
                "instruction_execution_icon": 10,
                "successful_auto_instructions": [
                  "test-instruction-correlation-alarm-get-6-2-name"
                ]
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-6-3"
                },
                "instruction_execution_icon": 3,
                "failed_auto_instructions": [
                  "test-instruction-correlation-alarm-get-6-3-name"
                ]
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-6-4"
                },
                "instruction_execution_icon": 11,
                "successful_manual_instructions": [
                  "test-instruction-correlation-alarm-get-6-4-name"
                ]
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-6-5"
                },
                "instruction_execution_icon": 4,
                "failed_manual_instructions": [
                  "test-instruction-correlation-alarm-get-6-5-name"
                ]
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
  Scenario: given get correlation request should return assigned_declare_ticket_rules for children
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-correlation-alarm-get-7",
        "connector_name": "test-connector-name-correlation-alarm-get-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-7",
        "resource": "test-resource-correlation-alarm-get-7-1",
        "state": 1,
        "output": "test-output-correlation-alarm-get-7"
      },
      {
        "connector": "test-connector-correlation-alarm-get-7",
        "connector_name": "test-connector-name-correlation-alarm-get-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-7",
        "resource": "test-resource-correlation-alarm-get-7-2",
        "state": 1,
        "output": "test-output-correlation-alarm-get-7"
      },
      {
        "connector": "test-connector-correlation-alarm-get-7",
        "connector_name": "test-connector-name-correlation-alarm-get-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-7",
        "resource": "test-resource-correlation-alarm-get-7-3",
        "state": 1,
        "output": "test-output-correlation-alarm-get-7"
      },
      {
        "connector": "test-connector-correlation-alarm-get-7",
        "connector_name": "test-connector-name-correlation-alarm-get-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-7",
        "resource": "test-resource-correlation-alarm-get-7-4",
        "state": 1,
        "output": "test-output-correlation-alarm-get-7"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-7&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response childAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response childAlarmID3={{ (index .lastResponse.data 2)._id }}
    When I save response childAlarmID4={{ (index .lastResponse.data 3)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metalarm-correlation-alarm-get-7-1",
      "comment": "test-metalarm-correlation-alarm-get-7-1-comment",
      "alarms": [
        "{{ .childAlarmID1 }}",
        "{{ .childAlarmID2 }}"
      ]
    }
    """
    Then the response code should be 204
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metalarm-correlation-alarm-get-7-2",
      "comment": "test-metalarm-correlation-alarm-get-7-2-comment",
      "alarms": [
        "{{ .childAlarmID3 }}",
        "{{ .childAlarmID4 }}"
      ]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-7&correlation=true&sort_by=v.display_name&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 2
        },
        {
          "children": 2
        }
      ]
    }
    """
    When I save response metaAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-alarm-correlation-rule-8",
      "system_name": "test-alarm-correlation-rule-8-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "GET",
            "auth": {
              "username": "test",
              "password": "test"
            },
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "declare_ticket": {
            "is_regexp": false,
            "ticket_id": "_id",
            "ticket_url": "url",
            "ticket_custom": "custom",
            "empty_response": false
          },
          "stop_on_fail": true
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-correlation-alarm-get-7-2"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-correlation-alarm-get-7-3"
            }
          }
        ]
      ]
    }
    """
    Then I save response ruleID={{ .lastResponse._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID1 }}",
        "with_declare_tickets": true,
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ .metaAlarmID2 }}",
        "with_declare_tickets": true,
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
                  "resource": "test-resource-correlation-alarm-get-7-1",
                  "state": {
                    "val": 1
                  }
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-7-2"
                },
                "assigned_declare_ticket_rules": [
                  {
                    "_id": "{{ .ruleID }}",
                    "name": "test-alarm-correlation-rule-8"
                  }
                ]
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
                  "resource": "test-resource-correlation-alarm-get-7-3"
                },
                "assigned_declare_ticket_rules": [
                  {
                    "_id": "{{ .ruleID }}",
                    "name": "test-alarm-correlation-rule-8"
                  }
                ]
              },
              {
                "v": {
                  "resource": "test-resource-correlation-alarm-get-7-4"
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
    Then the response key "0.data.children.data.0.assigned_declare_ticket_rules" should not exist
    Then the response key "1.data.children.data.1.assigned_declare_ticket_rules" should not exist

  @concurrent
  Scenario: given get search correlation request should return consistent opened/closed children values
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-correlation-alarm-get-8",
        "connector_name": "test-connector-name-correlation-alarm-get-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-8",
        "resource": "test-resource-correlation-alarm-get-8-1",
        "state": 1,
        "output": "test-output-correlation-alarm-get-8"
      },
      {
        "connector": "test-connector-correlation-alarm-get-8",
        "connector_name": "test-connector-name-correlation-alarm-get-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-8",
        "resource": "test-resource-correlation-alarm-get-8-2",
        "state": 1,
        "output": "test-resource-correlation-alarm-get-8-search"
      },
      {
        "connector": "test-connector-correlation-alarm-get-8",
        "connector_name": "test-connector-name-correlation-alarm-get-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-correlation-alarm-get-8",
        "resource": "test-resource-correlation-alarm-get-8-3",
        "state": 1,
        "output": "test-resource-correlation-alarm-get-8-search"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-8&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response childAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response childAlarmID3={{ (index .lastResponse.data 2)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metalarm-correlation-alarm-get-8-1",
      "comment": "test-metalarm-correlation-alarm-get-8-1-comment",
      "alarms": [
        "{{ .childAlarmID1 }}",
        "{{ .childAlarmID2 }}",
        "{{ .childAlarmID3 }}"
      ]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-8&correlation=true&sort_by=v.output&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 3,
          "opened_children": 3
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-correlation-alarm-get-8",
      "connector_name": "test-connector-name-correlation-alarm-get-8",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-correlation-alarm-get-8",
      "resource": "test-resource-correlation-alarm-get-8-3"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-correlation-alarm-get-8",
      "connector_name": "test-connector-name-correlation-alarm-get-8",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-correlation-alarm-get-8",
      "resource": "test-resource-correlation-alarm-get-8-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-8&correlation=true&sort_by=v.output&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 3,
          "opened_children": 2,
          "closed_children": 1
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-correlation-alarm-get-8",
      "connector_name": "test-connector-name-correlation-alarm-get-8",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-correlation-alarm-get-8",
      "resource": "test-resource-correlation-alarm-get-8-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-correlation-alarm-get-8",
      "connector_name": "test-connector-name-correlation-alarm-get-8",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-correlation-alarm-get-8",
      "resource": "test-resource-correlation-alarm-get-8-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-correlation-alarm-get-8",
      "connector_name": "test-connector-name-correlation-alarm-get-8",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-correlation-alarm-get-8",
      "resource": "test-resource-correlation-alarm-get-8-2"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-correlation-alarm-get-8",
      "connector_name": "test-connector-name-correlation-alarm-get-8",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-correlation-alarm-get-8",
      "resource": "test-resource-correlation-alarm-get-8-2"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-8&correlation=true&sort_by=v.output&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 3,
          "opened_children": 0,
          "closed_children": 3
        }
      ]
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-correlation-alarm-get-8",
      "connector_name": "test-connector-name-correlation-alarm-get-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-correlation-alarm-get-8",
      "resource": "test-resource-correlation-alarm-get-8-4",
      "state": 1,
      "output": "test-resource-correlation-alarm-get-8-search"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-8-4&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID4={{ (index .lastResponse.data 0)._id }}
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ .metaAlarmID }}/add:
    """json
    {
      "comment": "test-metalarm-correlation-alarm-get-8-1-comment",
      "alarms": ["{{ .childAlarmID4 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-8&correlation=true&sort_by=v.output&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 4,
          "opened_children": 1,
          "closed_children": 3
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ .metaAlarmID }}/remove:
    """json
    {
      "comment": "test-metalarm-correlation-alarm-get-8-1-comment",
      "alarms": ["{{ .childAlarmID4 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-8&correlation=true&sort_by=v.output&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 3,
          "opened_children": 0,
          "closed_children": 3
        },
        {}
      ]
    }
    """

  @concurrent
  Scenario: given meta alarm rule and alarms should increase events_count and last_event_date in metaalarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "_id": "test-correlation-alarm-get-9",
      "name": "test-correlation-alarm-get-9",
      "type": "complex",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-correlation-alarm-get-9"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 20,
          "unit": "s"
        },
        "threshold_count": 3
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-9",
      "connector_name": "test-correlation-alarm-get-9-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-9",
      "resource": "test-resource-correlation-alarm-get-9-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-9",
      "connector_name": "test-correlation-alarm-get-9-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-9",
      "resource": "test-resource-correlation-alarm-get-9-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-9",
      "connector_name": "test-correlation-alarm-get-9-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-9",
      "resource": "test-resource-correlation-alarm-get-9-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-9",
      "connector_name": "test-correlation-alarm-get-9-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-9",
      "resource": "test-resource-correlation-alarm-get-9-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-9-3
    Then the response code should be 200
    Then I save response lastEventDate={{ (index .lastResponse.data 0).v.last_event_date }}
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-9&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-correlation-alarm-get-9"
          },
          "v": {
            "events_count": 4
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
    Then the response key "data.0.v.last_event_date" should be "{{ .lastEventDate }}"
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-9",
      "connector_name": "test-correlation-alarm-get-9-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-9",
      "resource": "test-resource-correlation-alarm-get-9-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-9-3
    Then the response code should be 200
    Then I save response lastEventDate={{ (index .lastResponse.data 0).v.last_event_date }}
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-9&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-correlation-alarm-get-9"
          },
          "v": {
            "events_count": 5
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
    Then the response key "data.0.v.last_event_date" should be "{{ .lastEventDate }}"
    When I wait 2s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-9",
      "connector_name": "test-correlation-alarm-get-9-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-9",
      "resource": "test-resource-correlation-alarm-get-9-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-9-4
    Then the response code should be 200
    Then I save response lastEventDate={{ (index .lastResponse.data 0).v.last_event_date }}
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-9&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-correlation-alarm-get-9"
          },
          "v": {
            "events_count": 6
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
    Then the response key "data.0.v.last_event_date" should be "{{ .lastEventDate }}"
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-9",
      "connector_name": "test-correlation-alarm-get-9-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-9",
      "resource": "test-resource-correlation-alarm-get-9-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-9&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-correlation-alarm-get-9"
          },
          "v": {
            "events_count": 7
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "cancel",
      "connector": "test-correlation-alarm-get-9",
      "connector_name": "test-correlation-alarm-get-9-name",
      "source_type": "resource",
      "component":  "test-component-correlation-alarm-get-9",
      "resource": "test-resource-correlation-alarm-get-9-4",
      "output": "test"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "resolve_cancel",
      "connector": "test-correlation-alarm-get-9",
      "connector_name": "test-correlation-alarm-get-9-name",
      "source_type": "resource",
      "component":  "test-component-correlation-alarm-get-9",
      "resource": "test-resource-correlation-alarm-get-9-4"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-9-3
    Then the response code should be 200
    Then I save response lastEventDate={{ (index .lastResponse.data 0).v.last_event_date }}
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-9&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-correlation-alarm-get-9"
          },
          "v": {
            "events_count": 5
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
    Then the response key "data.0.v.last_event_date" should be "{{ .lastEventDate }}"
