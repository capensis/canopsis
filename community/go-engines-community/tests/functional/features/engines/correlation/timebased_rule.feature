Feature: correlation feature - timebased rule

  Scenario: given meta alarm rule and events should create meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-timebased-correlation-1",
      "type": "timebased",
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        }
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-timebased-correlation-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-timebased-1",
      "connector_name": "test-timebased-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-1",
      "resource": "test-timebased-correlation-resource-1-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """json
    {
      "connector": "test-timebased-1",
      "connector_name": "test-timebased-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-1",
      "resource": "test-timebased-correlation-resource-1-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-timebased-correlation-resource-1&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-timebased-correlation-1"
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
                  "connector": "test-timebased-1",
                  "connector_name": "test-timebased-1-name",
                  "component":  "test-timebased-correlation-1",
                  "resource": "test-timebased-correlation-resource-1-1"
                }
              },
              {
                "v": {
                  "connector": "test-timebased-1",
                  "connector_name": "test-timebased-1-name",
                  "component":  "test-timebased-correlation-1",
                  "resource": "test-timebased-correlation-resource-1-2"
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

  Scenario: given meta alarm rule and events should create 2 meta alarms because of 2 separate time intervals
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-timebased-correlation-2",
      "type": "timebased",
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        }
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-timebased-correlation-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-timebased-2",
      "connector_name": "test-timebased-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-2",
      "resource": "test-timebased-correlation-resource-2-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """json
    {
      "connector": "test-timebased-2",
      "connector_name": "test-timebased-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-2",
      "resource": "test-timebased-correlation-resource-2-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I wait 4s
    When I send an event:
    """json
    {
      "connector": "test-timebased-2",
      "connector_name": "test-timebased-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-2",
      "resource": "test-timebased-correlation-resource-2-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """json
    {
      "connector": "test-timebased-2",
      "connector_name": "test-timebased-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-2",
      "resource": "test-timebased-correlation-resource-2-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-timebased-2",
      "connector_name": "test-timebased-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-2",
      "resource": "test-timebased-correlation-resource-2-5",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-timebased-correlation-resource-2&correlation=true&sort_by=t&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-timebased-correlation-2"
          }
        },
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-timebased-correlation-2"
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
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ (index .lastResponse.data 1)._id }}",
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
                  "connector": "test-timebased-2",
                  "connector_name": "test-timebased-2-name",
                  "component":  "test-timebased-correlation-2",
                  "resource": "test-timebased-correlation-resource-2-3"
                }
              },
              {
                "v": {
                  "connector": "test-timebased-2",
                  "connector_name": "test-timebased-2-name",
                  "component":  "test-timebased-correlation-2",
                  "resource": "test-timebased-correlation-resource-2-4"
                }
              },
              {
                "v": {
                  "connector": "test-timebased-2",
                  "connector_name": "test-timebased-2-name",
                  "component":  "test-timebased-correlation-2",
                  "resource": "test-timebased-correlation-resource-2-5"
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-timebased-2",
                  "connector_name": "test-timebased-2-name",
                  "component":  "test-timebased-correlation-2",
                  "resource": "test-timebased-correlation-resource-2-1"
                }
              },
              {
                "v": {
                  "connector": "test-timebased-2",
                  "connector_name": "test-timebased-2-name",
                  "component":  "test-timebased-correlation-2",
                  "resource": "test-timebased-correlation-resource-2-2"
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

  Scenario: given meta alarm rule and events should create one single meta alarms because first group didn't reached default timebased threshold = 2 alarms
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-timebased-correlation-3",
      "type": "timebased",
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        }
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-timebased-correlation-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-timebased-3",
      "connector_name": "test-timebased-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-3",
      "resource": "test-timebased-correlation-resource-3-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 4s
    When I send an event:
    """json
    {
      "connector": "test-timebased-3",
      "connector_name": "test-timebased-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-3",
      "resource": "test-timebased-correlation-resource-3-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """json
    {
      "connector": "test-timebased-3",
      "connector_name": "test-timebased-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-3",
      "resource": "test-timebased-correlation-resource-3-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-timebased-correlation-resource-3&correlation=true&sort_by=v.meta&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-timebased-correlation-3"
          }
        },
        {
          "is_meta_alarm": false,
          "v": {
            "connector": "test-timebased-3",
            "connector_name": "test-timebased-3-name",
            "component":  "test-timebased-correlation-3",
            "resource": "test-timebased-correlation-resource-3-1"
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
                  "connector": "test-timebased-3",
                  "connector_name": "test-timebased-3-name",
                  "component":  "test-timebased-correlation-3",
                  "resource": "test-timebased-correlation-resource-3-2"
                }
              },
              {
                "v": {
                  "connector": "test-timebased-3",
                  "connector_name": "test-timebased-3-name",
                  "component":  "test-timebased-correlation-3",
                  "resource": "test-timebased-correlation-resource-3-3"
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

  Scenario: given deleted meta alarm rule should delete meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-timebased-correlation-4",
      "type": "timebased",
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        }
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-timebased-correlation-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-timebased-4",
      "connector_name": "test-timebased-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-4",
      "resource": "test-timebased-correlation-resource-4-1",
      "state": 2
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-timebased-4",
      "connector_name": "test-timebased-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-4",
      "resource": "test-timebased-correlation-resource-4-2",
      "state": 2
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-timebased-correlation-resource-4&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-timebased-correlation-4"
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
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do DELETE /api/v4/cat/metaalarmrules/{{ .metaAlarmRuleID }}
    Then the response code should be 204
    When I do GET /api/v4/alarms/{{ .metaAlarmID }}
    Then the response code should be 404
    When I do GET /api/v4/alarms?search=test-timebased-correlation-4&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-timebased-4",
            "connector_name": "test-timebased-4-name",
            "component":  "test-timebased-correlation-4",
            "resource": "test-timebased-correlation-resource-4-1",
            "parents": []
          }
        },
        {
          "v": {
            "connector": "test-timebased-4",
            "connector_name": "test-timebased-4-name",
            "component":  "test-timebased-correlation-4",
            "resource": "test-timebased-correlation-resource-4-2",
            "parents": []
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
