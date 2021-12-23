Feature: correlation feature - valuegroup rule with threshold rate
  Scenario: given meta alarm rule with threshold rate and events should create meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-valuegroup-correlation-rate-1",
      "type": "valuegroup",
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_rate": 0.4,
        "value_paths": [
          "entity.infos.valuegroupRate1.value"
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-1-connector",
      "connector_name": "test-valuegroup-rule-rate-1-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-1-component",
      "resource": "test-valuegroup-rule-rate-1-resource-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-1-connector",
      "connector_name": "test-valuegroup-rule-rate-1-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-1-component",
      "resource": "test-valuegroup-rule-rate-1-resource-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=t&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "consequences": {
            "total": 2
          },
          "metaalarm": true,
          "rule": {
            "name": "test-valuegroup-correlation-rate-1"
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

    Scenario: given meta alarm rule with threshold rate and events shouldn't create meta alarm because rate was recomputed by the new entity event after, metaalarm should be create after reaching the new rate
      Given I am admin
      When I do POST /api/v4/eventfilter/rules:
      """
      {
        "description" : "test-correlation-valuegroup-rate-2",
        "enabled": true,
        "patterns" : [{
            "connector" : "test-valuegroup-rule-rate-2-connector"
        }],
        "enabled" : true,
        "external_data" : {},
        "config": {
          "actions": [
            {
              "type" : "set_entity_info_from_template",
              "name" : "valuegroupRate2",
              "value" : "{{ `{{.Event.ExtraInfos.valuegroupRate2}}` }}",
              "description" : "valuegroupRate2"
            }
          ],
          "on_success": "pass",
          "on_failure": "pass"
        },
        "priority" : 10001,
        "type" : "enrichment"
      }
      """
      Then the response code should be 201
      When I do POST /api/v4/cat/metaalarmrules:
      """
      {
        "name": "test-valuegroup-correlation-rate-2",
        "type": "valuegroup",
        "config": {
          "time_interval": {
            "value": 10,
            "unit": "s"
          },
          "threshold_rate": 0.4,
          "value_paths": [
            "entity.infos.valuegroupRate2.value"
          ]
        }
      }
      """
      Then the response code should be 201
      Then I save response metaAlarmRuleID={{ .lastResponse._id }}
      When I wait the next periodical process
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-2-connector",
        "connector_name": "test-valuegroup-rule-rate-2-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-2-component",
        "resource": "test-valuegroup-rule-rate-2-resource-1",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 1 events processing
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-2-connector",
        "connector_name": "test-valuegroup-rule-rate-2-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-2-component",
        "resource": "test-valuegroup-rule-rate-2-resource-new",
        "valuegroupRate2": "1",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 1 events processing
      When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=t&sort_dir=asc
      Then the response code should be 200
      Then the response body should contain:
      """
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
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-2-connector",
        "connector_name": "test-valuegroup-rule-rate-2-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-2-component",
        "resource": "test-valuegroup-rule-rate-2-resource-2",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 2 events processing
      When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=t&sort_dir=asc
      Then the response code should be 200
      Then the response body should contain:
      """
      {
        "data": [
          {
            "consequences": {
              "total": 3
            },
            "metaalarm": true,
            "rule": {
              "name": "test-valuegroup-correlation-rate-2"
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

    Scenario: given meta alarm rule with threshold rate and events should create one single meta alarms because first group didn't reached threshold
      Given I am admin
      When I do POST /api/v4/cat/metaalarmrules:
      """
      {
        "name": "test-valuegroup-correlation-rate-3",
        "type": "valuegroup",
        "config": {
          "time_interval": {
            "value": 3,
            "unit": "s"
          },
          "threshold_rate": 0.4,
          "value_paths": [
            "entity.infos.valuegroupRate3.value"
          ]
        }
      }
      """
      Then the response code should be 201
      Then I save response metaAlarmRuleID={{ .lastResponse._id }}
      When I wait the next periodical process
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-3-connector",
        "connector_name": "test-valuegroup-rule-rate-3-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-3-component",
        "resource": "test-valuegroup-rule-rate-3-resource-1",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 1 events processing
      When I wait 4s
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-3-connector",
        "connector_name": "test-valuegroup-rule-rate-3-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-3-component",
        "resource": "test-valuegroup-rule-rate-3-resource-2",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 1 events processing
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-3-connector",
        "connector_name": "test-valuegroup-rule-rate-3-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-3-component",
        "resource": "test-valuegroup-rule-rate-3-resource-3",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 2 events processing
      When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=t&sort_dir=asc
      Then the response code should be 200
      Then the response body should contain:
      """
      {
        "data": [
          {
            "consequences": {
              "total": 2
            },
            "metaalarm": true,
            "rule": {
              "name": "test-valuegroup-correlation-rate-3"
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

    Scenario: given meta alarm rule with threshold rate and events should create 2 meta alarms because of 2 separate time intervals
      Given I am admin
      When I do POST /api/v4/cat/metaalarmrules:
      """
      {
        "name": "test-valuegroup-correlation-rate-4",
        "type": "valuegroup",
        "config": {
          "time_interval": {
            "value": 3,
            "unit": "s"
          },
          "threshold_rate": 0.4,
          "value_paths": [
            "entity.infos.valuegroupRate4.value"
          ]
        }
      }
      """
      Then the response code should be 201
      Then I save response metaAlarmRuleID={{ .lastResponse._id }}
      When I wait the next periodical process
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-4-connector",
        "connector_name": "test-valuegroup-rule-rate-4-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-4-component",
        "resource": "test-valuegroup-rule-rate-4-resource-1",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 1 events processing
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-4-connector",
        "connector_name": "test-valuegroup-rule-rate-4-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-4-component",
        "resource": "test-valuegroup-rule-rate-4-resource-2",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 2 events processing
      When I wait 4s
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-4-connector",
        "connector_name": "test-valuegroup-rule-rate-4-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-4-component",
        "resource": "test-valuegroup-rule-rate-4-resource-3",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 1 events processing
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-4-connector",
        "connector_name": "test-valuegroup-rule-rate-4-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-4-component",
        "resource": "test-valuegroup-rule-rate-4-resource-4",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 2 events processing
      When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=t&sort_dir=asc
      Then the response code should be 200
      Then the response body should contain:
      """
      {
        "data": [
          {
            "consequences": {
              "total": 2
            },
            "metaalarm": true,
            "rule": {
              "name": "test-valuegroup-correlation-rate-4"
            }
          },
          {
            "consequences": {
              "total": 2
            },
            "metaalarm": true,
            "rule": {
              "name": "test-valuegroup-correlation-rate-4"
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

    Scenario: given meta alarm rule with threshold rate and events should create one single meta alarm without first alarm, because interval shifting
      Given I am admin
      When I do POST /api/v4/cat/metaalarmrules:
      """
      {
        "name": "test-valuegroup-correlation-rate-5",
        "type": "valuegroup",
        "config": {
          "time_interval": {
            "value": 5,
            "unit": "s"
          },
          "threshold_rate": 0.6,
          "value_paths": [
            "entity.infos.valuegroupRate5.value"
          ]
        }
      }
      """
      Then the response code should be 201
      Then I save response metaAlarmRuleID={{ .lastResponse._id }}
      When I wait the next periodical process
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-5-connector",
        "connector_name": "test-valuegroup-rule-rate-5-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-5-component",
        "resource": "test-valuegroup-rule-rate-5-resource-1",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 1 events processing
      When I wait 3s
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-5-connector",
        "connector_name": "test-valuegroup-rule-rate-5-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-5-component",
        "resource": "test-valuegroup-rule-rate-5-resource-2",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 1 events processing
      When I wait 3s
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-5-connector",
        "connector_name": "test-valuegroup-rule-rate-5-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-5-component",
        "resource": "test-valuegroup-rule-rate-5-resource-3",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 1 events processing
      When I send an event:
      """
      {
        "connector": "test-valuegroup-rule-rate-5-connector",
        "connector_name": "test-valuegroup-rule-rate-5-connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-valuegroup-rule-rate-5-component",
        "resource": "test-valuegroup-rule-rate-5-resource-4",
        "state": 2,
        "output": "test",
        "long_output": "test",
        "author": "test-author"
      }
      """
      When I wait the end of 2 events processing
      When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=t&sort_dir=asc
      Then the response code should be 200
      Then the response body should contain:
      """
      {
        "data": [
          {
            "consequences": {
              "total": 3
            },
            "metaalarm": true,
            "rule": {
              "name": "test-valuegroup-correlation-rate-5"
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

  Scenario: given meta alarm rule with threshold rate and events should create meta alarm regarding total entity pattern
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-valuegroup-correlation-rate-6",
      "type": "valuegroup",
      "config": {
        "total_entity_patterns": [
          {
            "name": {
              "regex_match": "test-valuegroup-rule-rate-6-resource"
            }
          },
          {
            "name": {
              "regex_match": "test-valuegroup-rule-rate-7-resource"
            }
          }
        ],
        "time_interval": {
          "value": 15,
          "unit": "s"
        },
        "threshold_rate": 0.5,
        "value_paths": [
          "entity.infos.valuegroupRate6And7.value"
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-6-connector",
      "connector_name": "test-valuegroup-rule-rate-6-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-6-component",
      "resource": "test-valuegroup-rule-rate-6-resource-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-6-connector",
      "connector_name": "test-valuegroup-rule-rate-6-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-6-component",
      "resource": "test-valuegroup-rule-rate-6-resource-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-6-connector",
      "connector_name": "test-valuegroup-rule-rate-6-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-6-component",
      "resource": "test-valuegroup-rule-rate-6-resource-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-6-connector",
      "connector_name": "test-valuegroup-rule-rate-6-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-6-component",
      "resource": "test-valuegroup-rule-rate-6-resource-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=t&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
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
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-6-connector",
      "connector_name": "test-valuegroup-rule-rate-6-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-6-component",
      "resource": "test-valuegroup-rule-rate-6-resource-5",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=t&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "consequences": {
            "total": 5
          },
          "metaalarm": true,
          "rule": {
            "name": "test-valuegroup-correlation-rate-6"
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

  Scenario: given meta alarm rule with threshold rate and events should create 4 meta alarm regarding total counted by valuepath combinations
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-valuegroup-correlation-rate-8",
      "type": "valuegroup",
      "config": {
        "time_interval": {
          "value": 15,
          "unit": "s"
        },
        "threshold_rate": 0.5,
        "value_paths": [
          "entity.infos.valuegroupRate81.value",
          "entity.infos.valuegroupRate82.value"
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-8-connector",
      "connector_name": "test-valuegroup-rule-rate-8-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-8-component",
      "resource": "test-valuegroup-rule-rate-8-resource-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-8-connector",
      "connector_name": "test-valuegroup-rule-rate-8-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-8-component",
      "resource": "test-valuegroup-rule-rate-8-resource-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-8-connector",
      "connector_name": "test-valuegroup-rule-rate-8-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-8-component",
      "resource": "test-valuegroup-rule-rate-8-resource-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-8-connector",
      "connector_name": "test-valuegroup-rule-rate-8-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-8-component",
      "resource": "test-valuegroup-rule-rate-8-resource-5",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-8-connector",
      "connector_name": "test-valuegroup-rule-rate-8-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-8-component",
      "resource": "test-valuegroup-rule-rate-8-resource-7",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-8-connector",
      "connector_name": "test-valuegroup-rule-rate-8-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-8-component",
      "resource": "test-valuegroup-rule-rate-8-resource-8",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-8-connector",
      "connector_name": "test-valuegroup-rule-rate-8-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-8-component",
      "resource": "test-valuegroup-rule-rate-8-resource-10",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-8-connector",
      "connector_name": "test-valuegroup-rule-rate-8-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-8-component",
      "resource": "test-valuegroup-rule-rate-8-resource-11",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=t&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "consequences": {
            "total": 2
          },
          "metaalarm": true,
          "rule": {
            "name": "test-valuegroup-correlation-rate-8"
          }
        },
        {
          "consequences": {
            "total": 2
          },
          "metaalarm": true,
          "rule": {
            "name": "test-valuegroup-correlation-rate-8"
          }
        },
        {
          "consequences": {
            "total": 2
          },
          "metaalarm": true,
          "rule": {
            "name": "test-valuegroup-correlation-rate-8"
          }
        },
        {
          "consequences": {
            "total": 2
          },
          "metaalarm": true,
          "rule": {
            "name": "test-valuegroup-correlation-rate-8"
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
