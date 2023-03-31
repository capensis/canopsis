Feature: Get alarms
  I need to be able to get a alarms

  Scenario: given get correlation request should return meta alarms or alarms without parent
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-correlation-get-1",
        "connector_name": "test-connector-name-to-alarm-correlation-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-1",
        "resource": "test-resource-to-alarm-correlation-get-1-1",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-1"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-1",
        "connector_name": "test-connector-name-to-alarm-correlation-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-1",
        "resource": "test-resource-to-alarm-correlation-get-1-2",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-1"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-1",
        "connector_name": "test-connector-name-to-alarm-correlation-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-1",
        "resource": "test-resource-to-alarm-correlation-get-1-3",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-1"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-1",
        "connector_name": "test-connector-name-to-alarm-correlation-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-1",
        "resource": "test-resource-to-alarm-correlation-get-1-4",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-1"
      }
    ]
    """
    When I wait the end of 8 events processing
    When I do GET /api/v4/alarms?correlation=false&search=test-resource-to-alarm-correlation-get-1&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-1-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-1-2"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-1-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-1-4"
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
    When I do GET /api/v4/alarms?correlation=true&search=test-resource-to-alarm-correlation-get-1&multi_sort[]=v.component,asc&multi_sort[]=meta_alarm_rule.name,asc
    Then the response code should be 200
    Then the response body should contain:
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
            "output": "Rule: test-metaalarm-rule-alarm-correlation-get-1-1; Count: 3; Children: test-component-to-alarm-correlation-get-1"
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
            "output": "Rule: test-metaalarm-rule-alarm-correlation-get-1-2; Count: 1; Children: test-component-to-alarm-correlation-get-1",
            "children": [
              "test-resource-to-alarm-correlation-get-1-3/test-component-to-alarm-correlation-get-1"
            ]
          }
        },
        {
          "is_meta_alarm": false,
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-1-4"
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
                  "resource": "test-resource-to-alarm-correlation-get-1-1"
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
                  "resource": "test-resource-to-alarm-correlation-get-1-2"
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
                  "resource": "test-resource-to-alarm-correlation-get-1-3"
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
                  "resource": "test-resource-to-alarm-correlation-get-1-3"
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
                  "resource": "test-resource-to-alarm-correlation-get-1-1"
                }
              },
              {
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-1-2"
                }
              },
              {
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-1-3"
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

  Scenario: given get correlation request should return resolved children
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-correlation-get-2",
        "connector_name": "test-connector-name-to-alarm-correlation-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-2",
        "resource": "test-resource-to-alarm-correlation-get-2-1",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-2"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-2",
        "connector_name": "test-connector-name-to-alarm-correlation-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-2",
        "resource": "test-resource-to-alarm-correlation-get-2-2",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-2"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-2",
        "connector_name": "test-connector-name-to-alarm-correlation-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-2",
        "resource": "test-resource-to-alarm-correlation-get-2-3",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-2"
      }
    ]
    """
    When I wait the end of 8 events processing
    When I do GET /api/v4/alarms?correlation=false&search=test-resource-to-alarm-correlation-get-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response childAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response childAlarmID3={{ (index .lastResponse.data 2)._id }}
    When I do GET /api/v4/alarms?correlation=true&search=test-resource-to-alarm-correlation-get-2&multi_sort[]=v.component,asc&multi_sort[]=meta_alarm_rule.name,asc
    Then the response code should be 200
    When I save response metaAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response metaAlarmResource2={{ (index .lastResponse.data 1).v.resource }}
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-correlation-get-2",
        "connector_name": "test-connector-name-to-alarm-correlation-get-2",
        "source_type": "resource",
        "event_type": "cancel",
        "component": "test-component-to-alarm-correlation-get-2",
        "resource": "test-resource-to-alarm-correlation-get-2-2"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-2",
        "connector_name": "test-connector-name-to-alarm-correlation-get-2",
        "source_type": "resource",
        "event_type": "cancel",
        "component": "test-component-to-alarm-correlation-get-2",
        "resource": "test-resource-to-alarm-correlation-get-2-3"
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
    When I wait the end of 5 events processing
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-correlation-get-2",
        "connector_name": "test-connector-name-to-alarm-correlation-get-2",
        "source_type": "resource",
        "event_type": "resolve_cancel",
        "component": "test-component-to-alarm-correlation-get-2",
        "resource": "test-resource-to-alarm-correlation-get-2-2"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-2",
        "connector_name": "test-connector-name-to-alarm-correlation-get-2",
        "source_type": "resource",
        "event_type": "resolve_cancel",
        "component": "test-component-to-alarm-correlation-get-2",
        "resource": "test-resource-to-alarm-correlation-get-2-3"
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
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?correlation=true&opened=true&search=test-resource-to-alarm-correlation-get-2&multi_sort[]=v.component,asc&multi_sort[]=meta_alarm_rule.name,asc
    Then the response code should be 200
    Then the response body should contain:
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
    When I do GET /api/v4/alarms?correlation=true&opened=false&search=test-resource-to-alarm-correlation-get-2&multi_sort[]=v.component,asc&multi_sort[]=meta_alarm_rule.name,asc
    Then the response code should be 200
    Then the response body should contain:
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
                  "resource": "test-resource-to-alarm-correlation-get-2-1"
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
                  "resource": "test-resource-to-alarm-correlation-get-2-2"
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
                  "resource": "test-resource-to-alarm-correlation-get-2-3"
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
                  "resource": "test-resource-to-alarm-correlation-get-2-1"
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
                  "resource": "test-resource-to-alarm-correlation-get-2-1"
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
                  "resource": "test-resource-to-alarm-correlation-get-2-3"
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
                  "resource": "test-resource-to-alarm-correlation-get-2-3"
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

  Scenario: given get search correlation request should return filtered children
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-correlation-get-4",
        "connector_name": "test-connector-name-to-alarm-correlation-get-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-4",
        "resource": "test-resource-to-alarm-correlation-get-4-1",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-4"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-4",
        "connector_name": "test-connector-name-to-alarm-correlation-get-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-4",
        "resource": "test-resource-to-alarm-correlation-get-4-2",
        "state": 1,
        "output": "test-resource-to-alarm-correlation-get-4-search"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-4",
        "connector_name": "test-connector-name-to-alarm-correlation-get-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-4",
        "resource": "test-resource-to-alarm-correlation-get-4-3",
        "state": 1,
        "output": "test-resource-to-alarm-correlation-get-4-search"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-4",
        "connector_name": "test-connector-name-to-alarm-correlation-get-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-4",
        "resource": "test-resource-to-alarm-correlation-get-4-4",
        "state": 1,
        "output": "test-resource-to-alarm-correlation-get-4-search"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-4",
        "connector_name": "test-connector-name-to-alarm-correlation-get-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-4",
        "resource": "test-resource-to-alarm-correlation-get-4-5",
        "state": 1,
        "output": "test-resource-to-alarm-correlation-get-4-search"
      }
    ]
    """
    When I wait the end of 5 events processing
    When I send an event:
    """json
    {
      "connector": "engine",
      "connector_name": "correlation",
      "source_type": "component",
      "event_type": "manual_metaalarm_group",
      "component":  "metaalarm",
      "output": "test-output-to-alarm-correlation-get-4",
      "display_name": "test-metalarm-to-alarm-correlation-get-4-1",
      "ma_children": [
        "test-resource-to-alarm-correlation-get-4-1/test-component-to-alarm-correlation-get-4",
        "test-resource-to-alarm-correlation-get-4-2/test-component-to-alarm-correlation-get-4",
        "test-resource-to-alarm-correlation-get-4-3/test-component-to-alarm-correlation-get-4"
      ]
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "engine",
      "connector_name": "correlation",
      "source_type": "component",
      "event_type": "manual_metaalarm_group",
      "component":  "metaalarm",
      "output": "test-output-to-alarm-correlation-get-4",
      "display_name": "test-metalarm-to-alarm-correlation-get-4-2",
      "ma_children": [
        "test-resource-to-alarm-correlation-get-4-4/test-component-to-alarm-correlation-get-4"
      ]
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-correlation-get-4",
      "connector_name": "test-connector-name-to-alarm-correlation-get-4",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-to-alarm-correlation-get-4",
      "resource": "test-resource-to-alarm-correlation-get-4-3"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-correlation-get-4",
      "connector_name": "test-connector-name-to-alarm-correlation-get-4",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-to-alarm-correlation-get-4",
      "resource": "test-resource-to-alarm-correlation-get-4-3"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?correlation=false&search=test-resource-to-alarm-correlation-get-4&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response childAlarmID3={{ (index .lastResponse.data 2)._id }}
    When I save response childAlarmID4={{ (index .lastResponse.data 3)._id }}
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-correlation-get-4-search&active_columns[]=v.output&correlation=true&sort_by=children&sort=desc
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
          "children": 1,
          "filtered_children": [
            "{{ .childAlarmID4 }}"
          ]
        },
        {
          "is_meta_alarm": false,
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-4-5"
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
    Then the response array key "data.0.filtered_children" should contain:
    """json
    [
      "{{ .childAlarmID2 }}",
      "{{ .childAlarmID3 }}"
    ]
    """
    When I do GET /api/v4/alarms?opened=true&search=test-resource-to-alarm-correlation-get-4-search&active_columns[]=v.output&correlation=true&sort_by=children&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "filtered_children": [
            "{{ .childAlarmID2 }}"
          ]
        },
        {
          "is_meta_alarm": true,
          "children": 1,
          "filtered_children": [
            "{{ .childAlarmID4 }}"
          ]
        },
        {
          "is_meta_alarm": false,
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-4-5"
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

  Scenario: given get correlation with_instructions request should return
  meta alarms with children, children should have assigned instruction if they have it, the corresponding
  metaalarm should have a mark about it
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-correlation-get-5",
        "connector_name": "test-connector-name-to-alarm-correlation-get-5",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-5",
        "resource": "test-resource-to-alarm-correlation-get-5-1",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-5"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-5",
        "connector_name": "test-connector-name-to-alarm-correlation-get-5",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-5",
        "resource": "test-resource-to-alarm-correlation-get-5-2",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-5"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-5",
        "connector_name": "test-connector-name-to-alarm-correlation-get-5",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-5",
        "resource": "test-resource-to-alarm-correlation-get-5-3",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-5"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-5",
        "connector_name": "test-connector-name-to-alarm-correlation-get-5",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-5",
        "resource": "test-resource-to-alarm-correlation-get-5-4",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-5"
      }
    ]
    """
    When I wait the end of 4 events processing
    When I send an event:
    """json
    [
      {
        "connector": "engine",
        "connector_name": "correlation",
        "source_type": "component",
        "event_type": "manual_metaalarm_group",
        "component":  "metaalarm",
        "output": "test-output-to-alarm-correlation-get-5",
        "display_name": "test-metalarm-to-alarm-correlation-get-5-1",
        "ma_children": [
          "test-resource-to-alarm-correlation-get-5-1/test-component-to-alarm-correlation-get-5",
          "test-resource-to-alarm-correlation-get-5-2/test-component-to-alarm-correlation-get-5"
        ]
      },
      {
        "connector": "engine",
        "connector_name": "correlation",
        "source_type": "component",
        "event_type": "manual_metaalarm_group",
        "component":  "metaalarm",
        "output": "test-output-to-alarm-correlation-get-5",
        "display_name": "test-metalarm-to-alarm-correlation-get-5-2",
        "ma_children": [
          "test-resource-to-alarm-correlation-get-5-3/test-component-to-alarm-correlation-get-5",
          "test-resource-to-alarm-correlation-get-5-4/test-component-to-alarm-correlation-get-5"
        ]
      }
    ]
    """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-correlation-get-5&correlation=true&with_instructions=true&sort_by=v.display_name&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children_instructions": true,
          "v": {
            "display_name": "test-metalarm-to-alarm-correlation-get-5-1"
          }
        },
        {
          "is_meta_alarm": true,
          "children_instructions": false,
          "v": {
            "display_name": "test-metalarm-to-alarm-correlation-get-5-2"
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
                  "resource": "test-resource-to-alarm-correlation-get-5-1",
                  "state": {
                    "val": 1
                  }
                },
                "assigned_instructions": []
              },
              {
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-5-2"
                },
                "assigned_instructions": [
                  {
                    "_id": "test-instruction-to-alarm-correlation-instruction-get-1-1"
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
                  "resource": "test-resource-to-alarm-correlation-get-5-3"
                },
                "assigned_instructions": []
              },
              {
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-5-4"
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

  Scenario: given get correlation alarms with the same children, but different alarms shouldn't have alarms of each other
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-correlation-get-6",
        "connector_name": "test-connector-name-to-alarm-correlation-get-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-6",
        "resource": "test-resource-to-alarm-correlation-get-6-1",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-6"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-6",
        "connector_name": "test-connector-name-to-alarm-correlation-get-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-6",
        "resource": "test-resource-to-alarm-correlation-get-6-2",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-6"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-correlation-get-6&opened=true&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response childAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I send an event:
    """json
    {
      "connector": "engine",
      "connector_name": "correlation",
      "source_type": "component",
      "event_type": "manual_metaalarm_group",
      "component":  "metaalarm",
      "output": "test-output-to-alarm-correlation-get-6",
      "display_name": "test-metalarm-to-alarm-correlation-get-6-1",
      "ma_children": [
        "test-resource-to-alarm-correlation-get-6-1/test-component-to-alarm-correlation-get-6",
        "test-resource-to-alarm-correlation-get-6-2/test-component-to-alarm-correlation-get-6"
      ]
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-correlation-get-6",
      "connector_name": "test-connector-name-to-alarm-correlation-get-6",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-to-alarm-correlation-get-6",
      "resource": "test-resource-to-alarm-correlation-get-6-1",
      "output": "test-output-to-alarm-correlation-get-6"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-correlation-get-6",
      "connector_name": "test-connector-name-to-alarm-correlation-get-6",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-to-alarm-correlation-get-6",
      "resource": "test-resource-to-alarm-correlation-get-6-1",
      "output": "test-output-to-alarm-correlation-get-6"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-correlation-get-6",
      "connector_name": "test-connector-name-to-alarm-correlation-get-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-correlation-get-6",
      "resource": "test-resource-to-alarm-correlation-get-6-1",
      "state": 2,
      "output": "test-output-to-alarm-correlation-get-6"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "engine",
      "connector_name": "correlation",
      "source_type": "component",
      "event_type": "manual_metaalarm_group",
      "component":  "metaalarm",
      "output": "test-output-to-alarm-correlation-get-6",
      "display_name": "test-metalarm-to-alarm-correlation-get-6-2",
      "ma_children": [
        "test-resource-to-alarm-correlation-get-6-1/test-component-to-alarm-correlation-get-6",
        "test-resource-to-alarm-correlation-get-6-2/test-component-to-alarm-correlation-get-6"
      ]
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-correlation-get-6&opened=true&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID3={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-correlation-get-6&correlation=true&sort_by=v.display_name&sort=asc
    Then the response code should be 200
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
                  "resource": "test-resource-to-alarm-correlation-get-6-1",
                  "state": {
                    "val": 1
                  }
                }
              },
              {
                "_id": "{{ .childAlarmID2 }}",
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-6-2"
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
                  "resource": "test-resource-to-alarm-correlation-get-6-1",
                  "state": {
                    "val": 2
                  }
                }
              },
              {
                "_id": "{{ .childAlarmID2 }}",
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-6-2"
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

  Scenario: given alarms which are matched to metaalarm should return instruction statuses and icons
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-correlation-get-7",
        "connector_name": "test-connector-name-to-alarm-correlation-get-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-7",
        "resource": "test-resource-to-alarm-correlation-get-7-1",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-7"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-7",
        "connector_name": "test-connector-name-to-alarm-correlation-get-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-7",
        "resource": "test-resource-to-alarm-correlation-get-7-2",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-7"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-7",
        "connector_name": "test-connector-name-to-alarm-correlation-get-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-7",
        "resource": "test-resource-to-alarm-correlation-get-7-3",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-7"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-7",
        "connector_name": "test-connector-name-to-alarm-correlation-get-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-7",
        "resource": "test-resource-to-alarm-correlation-get-7-4",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-7"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-7",
        "connector_name": "test-connector-name-to-alarm-correlation-get-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-7",
        "resource": "test-resource-to-alarm-correlation-get-7-5",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-7"
      }
    ]
    """
    When I wait the end of 14 events processing
    When I do GET /api/v4/alarms?correlation=false&search=test-resource-to-alarm-correlation-get-7&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-7-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-7-2"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-7-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-7-4"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-7-5"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 5
      }
    }
    """
    When I save response childAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response childAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response childAlarmID3={{ (index .lastResponse.data 2)._id }}
    When I save response childAlarmID4={{ (index .lastResponse.data 3)._id }}
    When I save response childAlarmID5={{ (index .lastResponse.data 4)._id }}
    When I do GET /api/v4/alarms?correlation=true&search=test-resource-to-alarm-correlation-get-7&multi_sort[]=v.component,asc&multi_sort[]=meta_alarm_rule.name,asc
    Then the response code should be 200
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .childAlarmID4 }}",
      "instruction": "test-instruction-to-alarm-correlation-get-7-4"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .childAlarmID5 }}",
      "instruction": "test-instruction-to-alarm-correlation-get-7-5"
    }
    """
    Then the response code should be 200
    When I wait the end of 4 events processing
    When I wait 3s
    When I do POST /api/v4/alarm-details:
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
                  "resource": "test-resource-to-alarm-correlation-get-7-1"
                },
                "instruction_execution_icon": 9
              },
              {
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-7-2"
                },
                "instruction_execution_icon": 10,
                "successful_auto_instructions": [
                  "test-instruction-to-alarm-correlation-get-7-2-name"
                ]
              },
              {
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-7-3"
                },
                "instruction_execution_icon": 3,
                "failed_auto_instructions": [
                  "test-instruction-to-alarm-correlation-get-7-3-name"
                ]
              },
              {
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-7-4"
                },
                "instruction_execution_icon": 11,
                "successful_manual_instructions": [
                  "test-instruction-to-alarm-correlation-get-7-4-name"
                ]
              },
              {
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-7-5"
                },
                "instruction_execution_icon": 4,
                "failed_manual_instructions": [
                  "test-instruction-to-alarm-correlation-get-7-5-name"
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

  Scenario: given get correlation request should return assigned_declare_ticket_rules for children
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-correlation-get-8",
        "connector_name": "test-connector-name-to-alarm-correlation-get-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-8",
        "resource": "test-resource-to-alarm-correlation-get-8-1",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-8"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-8",
        "connector_name": "test-connector-name-to-alarm-correlation-get-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-8",
        "resource": "test-resource-to-alarm-correlation-get-8-2",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-8"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-8",
        "connector_name": "test-connector-name-to-alarm-correlation-get-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-8",
        "resource": "test-resource-to-alarm-correlation-get-8-3",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-8"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-8",
        "connector_name": "test-connector-name-to-alarm-correlation-get-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-8",
        "resource": "test-resource-to-alarm-correlation-get-8-4",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-8"
      }
    ]
    """
    When I wait the end of 4 events processing
    When I send an event:
    """json
    [
      {
        "connector": "engine",
        "connector_name": "correlation",
        "source_type": "component",
        "event_type": "manual_metaalarm_group",
        "component":  "metaalarm",
        "output": "test-output-to-alarm-correlation-get-8",
        "display_name": "test-metalarm-to-alarm-correlation-get-8-1",
        "ma_children": [
          "test-resource-to-alarm-correlation-get-8-1/test-component-to-alarm-correlation-get-8",
          "test-resource-to-alarm-correlation-get-8-2/test-component-to-alarm-correlation-get-8"
        ]
      },
      {
        "connector": "engine",
        "connector_name": "correlation",
        "source_type": "component",
        "event_type": "manual_metaalarm_group",
        "component":  "metaalarm",
        "output": "test-output-to-alarm-correlation-get-8",
        "display_name": "test-metalarm-to-alarm-correlation-get-8-2",
        "ma_children": [
          "test-resource-to-alarm-correlation-get-8-3/test-component-to-alarm-correlation-get-8",
          "test-resource-to-alarm-correlation-get-8-4/test-component-to-alarm-correlation-get-8"
        ]
      }
    ]
    """
    When I wait the end of 4 events processing
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
              "value": "test-resource-to-alarm-correlation-get-8-2"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-alarm-correlation-get-8-3"
            }
          }
        ]
      ]
    }
    """
    Then I save response ruleID={{ .lastResponse._id }}
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-correlation-get-8&correlation=true&with_instructions=true&sort_by=v.display_name&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "v": {
            "display_name": "test-metalarm-to-alarm-correlation-get-8-1"
          }
        },
        {
          "is_meta_alarm": true,
          "v": {
            "display_name": "test-metalarm-to-alarm-correlation-get-8-2"
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
        "with_declare_tickets": true,
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ (index .lastResponse.data 1)._id }}",
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
                  "resource": "test-resource-to-alarm-correlation-get-8-1",
                  "state": {
                    "val": 1
                  }
                }
              },
              {
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-8-2"
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
                  "resource": "test-resource-to-alarm-correlation-get-8-3"
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
                  "resource": "test-resource-to-alarm-correlation-get-8-4"
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
