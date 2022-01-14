Feature: add infos to existed alarm
  I need to be able to add infos to existed alarm and later update rule to remove just added infos

  Scenario: given new dynamic infos and check update existed alarm
    Given I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-upd-existed-alarm",
      "description": "test-dynamicinfos-upd-existed-alarm-description",
      "disable_during_periods": [],
      "alarm_patterns": [
        {
          "v": {
            "output": {
                "regex_match": "dynamic-info-update alarm"
            },
            "component": "dynamic-info-update-component"
          }
        }
      ],
      "infos": [
        {
          "name": "test-dynamicinfos-upd-existed-alarm-name-1",
          "value": "Alarm field: {{ `{{ .Alarm.Value.Connector }}` }}; Entity field: {{ `{{ .Entity.ID }}` }};"
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
    When I wait the next periodical process
    When I do GET /api/v4/alarms?filter={"v.resource":"dynamic-info-update-resource-1"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
        "infos": {
            "{{ .ruleID }}": {
            "test-dynamicinfos-upd-existed-alarm-name-1": "Alarm field: dynamic-info-update-connector; Entity field: dynamic-info-update-resource-1/dynamic-info-update-component;"
            }
        },
          "v": {
            "component": "dynamic-info-update-component",
            "connector": "dynamic-info-update-connector",
            "connector_name": "dynamic-info-update-connectorname",
            "infos": {
              "{{ .ruleID }}": {
                "test-dynamicinfos-upd-existed-alarm-name-1": "Alarm field: dynamic-info-update-connector; Entity field: dynamic-info-update-resource-1/dynamic-info-update-component;"
              }
            },
            "resource": "dynamic-info-update-resource-1"
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
    When I do PUT /api/v4/cat/dynamic-infos/{{ .ruleID }}:
    """json
    {
      "name": "test-dynamicinfos-upd-existed-alarm",
      "description": "test-dynamicinfos-upd-existed-alarm-description",
      "disable_during_periods": [],
      "alarm_patterns": [
        {
          "v": {
            "output": {
                "regex_match": "dynamic-info-no-update alarm"
            },
            "component": "dynamic-info-no-update-component"
          }
        }
      ],
      "infos": [
        {
          "name": "test-dynamicinfos-upd-existed-alarm-name-1",
          "value": "Alarm field: {{ `{{ .Alarm.Value.Connector }}` }}; Entity field: {{ `{{ .Entity.ID }}` }};"
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
    When I wait the next periodical process
    When I do GET /api/v4/alarms?filter={"v.infos.{{ .ruleID }}":{"$exists":true}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  Scenario: given new dynamic infos with ".Event" reference and check it shouldn't update any existed alarm
    Given I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-shouldnt-upd-existed-alarm",
      "description": "test-dynamicinfos-shouldnt-upd-existed-alarm-description",
      "disable_during_periods": [],
      "alarm_patterns": [
        {
          "v": {
            "output": {
                "regex_match": "dynamic-info-update alarm"
            },
            "component": "dynamic-info-update-component"
          }
        }
      ],
      "infos": [
        {
          "name": "test-dynamicinfos-upd-existed-alarm-name-2",
          "value": "Alarm field: {{ `{{ .Alarm.Value.Connector }}` }}; Entity field: {{ `{{ .Entity.ID }}` }}; Event field: {{ `{{ .Event.ExtraInfos.customer }}` }}"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    Then the response key "alarm_update" should not exist
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I do GET /api/v4/alarms?filter={"v.infos.{{ .ruleID }}":{"$exists":true}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """