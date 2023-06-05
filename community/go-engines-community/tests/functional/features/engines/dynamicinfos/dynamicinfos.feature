Feature: add infos to alarm
  I need to be able to add infos to alarm

  Scenario: given dynamic infos should update new alarm
    Given I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-1-name",
      "description": "test-dynamicinfos-dynamicinfos-1-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamicinfos-1"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamicinfos-infos-1-name",
          "value": "Alarm field: {{ `{{ .Alarm.Value.Connector }}` }}; Entity field: {{ `{{ .Entity.ID }}` }}; Event field: {{ `{{ .Event.ExtraInfos.customer }}` }}"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "alarm_update": false
    }
    """
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector" : "test-connector-dynamicinfos-1",
      "connector_name" : "test-connector-name-dynamicinfos-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-dynamicinfos-1",
      "resource" : "test-resource-dynamicinfos-1",
      "state" : 1,
      "output" : "test-output-dynamicinfos-1",
      "customer": "test-customer-dynamicinfos-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filters[]=test-widgetfilter-dynamic-infos-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "infos": {
            "{{ .ruleID }}": {
              "test-dynamicinfos-infos-1-name": "Alarm field: test-connector-dynamicinfos-1; Entity field: test-resource-dynamicinfos-1/test-component-dynamicinfos-1; Event field: test-customer-dynamicinfos-1"
            }
          },
          "v": {
            "component": "test-component-dynamicinfos-1",
            "connector": "test-connector-dynamicinfos-1",
            "connector_name": "test-connector-name-dynamicinfos-1",
            "resource": "test-resource-dynamicinfos-1",
            "infos": {
              "{{ .ruleID }}": {
                "test-dynamicinfos-infos-1-name": "Alarm field: test-connector-dynamicinfos-1; Entity field: test-resource-dynamicinfos-1/test-component-dynamicinfos-1; Event field: test-customer-dynamicinfos-1"
              }
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

  Scenario: given new dynamic infos should update existed alarm
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-dynamicinfos-2",
      "connector_name" : "test-connector-name-dynamicinfos-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-dynamicinfos-2",
      "resource" : "test-resource-dynamicinfos-2",
      "state" : 1,
      "output" : "test-output-dynamicinfos-2",
      "customer": "test-customer-dynamicinfos-2"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-2-name",
      "description": "test-dynamicinfos-dynamicinfos-2-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamicinfos-2"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamicinfos-infos-2-name",
          "value": "Alarm field: {{ `{{ .Alarm.Value.Connector }}` }}; Entity field: {{ `{{ .Entity.ID }}` }}"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "alarm_update": true
    }
    """
    When I save response ruleID={{ .lastResponse._id }}
    When I do GET /api/v4/alarms?search=test-resource-dynamicinfos-2 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "infos": {
            "{{ .ruleID }}": {
              "test-dynamicinfos-infos-2-name": "Alarm field: test-connector-dynamicinfos-2; Entity field: test-resource-dynamicinfos-2/test-component-dynamicinfos-2"
            }
          },
          "v": {
            "component": "test-component-dynamicinfos-2",
            "connector": "test-connector-dynamicinfos-2",
            "connector_name": "test-connector-name-dynamicinfos-2",
            "resource": "test-resource-dynamicinfos-2",
            "infos": {
              "{{ .ruleID }}": {
                "test-dynamicinfos-infos-2-name": "Alarm field: test-connector-dynamicinfos-2; Entity field: test-resource-dynamicinfos-2/test-component-dynamicinfos-2"
              }
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

  Scenario: given new dynamic infos with event in template should not update new alarm
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-dynamicinfos-3",
      "connector_name" : "test-connector-name-dynamicinfos-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-dynamicinfos-3",
      "resource" : "test-resource-dynamicinfos-3",
      "state" : 1,
      "output" : "test-output-dynamicinfos-3",
      "customer": "test-customer-dynamicinfos-3"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-3-name",
      "description": "test-dynamicinfos-dynamicinfos-3-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamicinfos-3"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamicinfos-infos-3-name",
          "value": "Alarm field: {{ `{{ .Alarm.Value.Connector }}` }}; Entity field: {{ `{{ .Entity.ID }}` }}; Event field: {{ `{{ .Event.ExtraInfos.customer }}` }}"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "alarm_update": false
    }
    """
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I do GET /api/v4/alarms?filters[]=test-widgetfilter-dynamic-infos-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-dynamicinfos-3",
            "connector": "test-connector-dynamicinfos-3",
            "connector_name": "test-connector-name-dynamicinfos-3",
            "resource": "test-resource-dynamicinfos-3"
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

  Scenario: given updated dynamic infos should update infos in alarm
    Given I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-4-name",
      "description": "test-dynamicinfos-dynamicinfos-4-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamicinfos-4"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamicinfos-infos-4-name",
          "value": "Alarm field: {{ `{{ .Alarm.Value.Connector }}` }}; Entity field: {{ `{{ .Entity.ID }}` }}"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector" : "test-connector-dynamicinfos-4",
      "connector_name" : "test-connector-name-dynamicinfos-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-dynamicinfos-4",
      "resource" : "test-resource-dynamicinfos-4",
      "state" : 1,
      "output" : "test-output-dynamicinfos-4",
      "customer": "test-customer-dynamicinfos-4"
    }
    """
    When I wait the end of event processing
    When I do PUT /api/v4/cat/dynamic-infos/{{ .ruleID }}:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-4-name",
      "description": "test-dynamicinfos-dynamicinfos-4-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamicinfos-4"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamicinfos-infos-4-name",
          "value": "Alarm field updated: {{ `{{ .Alarm.Value.Component }}` }}; Entity field updated: {{ `{{ .Entity.Name }}` }}"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "alarm_update": true
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-dynamicinfos-4 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "infos": {
            "{{ .ruleID }}": {
              "test-dynamicinfos-infos-4-name": "Alarm field updated: test-component-dynamicinfos-4; Entity field updated: test-resource-dynamicinfos-4"
            }
          },
          "v": {
            "component": "test-component-dynamicinfos-4",
            "connector": "test-connector-dynamicinfos-4",
            "connector_name": "test-connector-name-dynamicinfos-4",
            "resource": "test-resource-dynamicinfos-4",
            "infos": {
              "{{ .ruleID }}": {
                "test-dynamicinfos-infos-4-name": "Alarm field updated: test-component-dynamicinfos-4; Entity field updated: test-resource-dynamicinfos-4"
              }
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

  Scenario: given updated dynamic infos should remove infos from alarm
    Given I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-5-name",
      "description": "test-dynamicinfos-dynamicinfos-5-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamicinfos-5"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamicinfos-infos-5-name",
          "value": "Alarm field: {{ `{{ .Alarm.Value.Connector }}` }}; Entity field: {{ `{{ .Entity.ID }}` }}; Event field: {{ `{{ .Event.ExtraInfos.customer }}` }}"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector" : "test-connector-dynamicinfos-5",
      "connector_name" : "test-connector-name-dynamicinfos-5",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-dynamicinfos-5",
      "resource" : "test-resource-dynamicinfos-5",
      "state" : 1,
      "output" : "test-output-dynamicinfos-5",
      "customer": "test-customer-dynamicinfos-5"
    }
    """
    When I wait the end of event processing
    When I do PUT /api/v4/cat/dynamic-infos/{{ .ruleID }}:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-5-name",
      "description": "test-dynamicinfos-dynamicinfos-5-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamicinfos-5-not-exist"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamicinfos-infos-5-name",
          "value": "Alarm field: {{ `{{ .Alarm.Value.Connector }}` }}; Entity field: {{ `{{ .Entity.ID }}` }}; Event field: {{ `{{ .Event.ExtraInfos.customer }}` }}"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "alarm_update": true
    }
    """
    When I do GET /api/v4/alarms?filters[]=test-widgetfilter-dynamic-infos-5 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-dynamicinfos-5",
            "connector": "test-connector-dynamicinfos-5",
            "connector_name": "test-connector-name-dynamicinfos-5",
            "resource": "test-resource-dynamicinfos-5"
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

  Scenario: given removed dynamic infos should remove infos from alarm
    Given I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-6-name",
      "description": "test-dynamicinfos-dynamicinfos-6-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamicinfos-6"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamicinfos-infos-6-name",
          "value": "Alarm field: {{ `{{ .Alarm.Value.Connector }}` }}; Entity field: {{ `{{ .Entity.ID }}` }}; Event field: {{ `{{ .Event.ExtraInfos.customer }}` }}"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector" : "test-connector-dynamicinfos-6",
      "connector_name" : "test-connector-name-dynamicinfos-6",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-dynamicinfos-6",
      "resource" : "test-resource-dynamicinfos-6",
      "state" : 1,
      "output" : "test-output-dynamicinfos-6",
      "customer": "test-customer-dynamicinfos-6"
    }
    """
    When I wait the end of event processing
    When I do DELETE /api/v4/cat/dynamic-infos/{{ .ruleID }}
    Then the response code should be 204
    When I do GET /api/v4/alarms?filters[]=test-widgetfilter-dynamic-infos-6 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-dynamicinfos-6",
            "connector": "test-connector-dynamicinfos-6",
            "connector_name": "test-connector-name-dynamicinfos-6",
            "resource": "test-resource-dynamicinfos-6"
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

  Scenario: given disabled dynamic infos should remove infos from alarm
    Given I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-7-name",
      "description": "test-dynamicinfos-dynamicinfos-7-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamicinfos-7"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamicinfos-infos-7-name",
          "value": "Alarm field: {{ `{{ .Alarm.Value.Connector }}` }}; Entity field: {{ `{{ .Entity.ID }}` }}; Event field: {{ `{{ .Event.ExtraInfos.customer }}` }}"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector" : "test-connector-dynamicinfos-7",
      "connector_name" : "test-connector-name-dynamicinfos-7",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-dynamicinfos-7",
      "resource" : "test-resource-dynamicinfos-7",
      "state" : 1,
      "output" : "test-output-dynamicinfos-7",
      "customer": "test-customer-dynamicinfos-7"
    }
    """
    When I wait the end of event processing
    When I do PUT /api/v4/cat/dynamic-infos/{{ .ruleID }}:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-7-name",
      "description": "test-dynamicinfos-dynamicinfos-7-description",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-dynamicinfos-7"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamicinfos-infos-7-name",
          "value": "Alarm field: {{ `{{ .Alarm.Value.Connector }}` }}; Entity field: {{ `{{ .Entity.ID }}` }}; Event field: {{ `{{ .Event.ExtraInfos.customer }}` }}"
        }
      ],
      "enabled": false
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "alarm_update": true
    }
    """
    When I do GET /api/v4/alarms?filters[]=test-widgetfilter-dynamic-infos-7 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-dynamicinfos-7",
            "connector": "test-connector-dynamicinfos-7",
            "connector_name": "test-connector-name-dynamicinfos-7",
            "resource": "test-resource-dynamicinfos-7"
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
