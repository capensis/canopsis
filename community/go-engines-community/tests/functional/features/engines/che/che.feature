Feature: create entities on event
  I need to be able to create entities on event

  Scenario: given resource check event should create entities
    Given I am admin
    When I send an event:
    """
    {
      "connector": "test-connector-che-1",
      "connector_name": "test-connector-name-che-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-1",
      "resource": "test-resource-che-1",
      "state": 2,
      "output": "test-output-che-1"
    }
    """
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=che-1
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-component-che-1",
          "category": null,
          "component": "test-component-che-1",
          "depends": [
            "test-resource-che-1/test-component-che-1"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-connector-che-1/test-connector-name-che-1"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-component-che-1",
          "type": "component"
        },
        {
          "_id": "test-connector-che-1/test-connector-name-che-1",
          "category": null,
          "depends": [
            "test-component-che-1"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-resource-che-1/test-component-che-1"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-connector-name-che-1",
          "type": "connector"
        },
        {
          "_id": "test-resource-che-1/test-component-che-1",
          "category": null,
          "component": "test-component-che-1",
          "depends": [
            "test-connector-che-1/test-connector-name-che-1"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-component-che-1"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-resource-che-1",
          "type": "resource"
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

  Scenario: given component event should create entities
    Given I am admin
    When I send an event:
    """
    {
      "connector": "test-connector-che-2",
      "connector_name": "test-connector-name-che-2",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-2",
      "state": 2,
      "output": "test-output-che-2"
    }
    """
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=che-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-component-che-2",
          "category": null,
          "component": "test-component-che-2",
          "depends": [],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-connector-che-2/test-connector-name-che-2"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-component-che-2",
          "type": "component"
        },
        {
          "_id": "test-connector-che-2/test-connector-name-che-2",
          "category": null,
          "depends": [
            "test-component-che-2"
          ],
          "enable_history": [
            {{ .createTimestamp }}
          ],
          "enabled": true,
          "impact": [],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-connector-name-che-2",
          "type": "connector"
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

  Scenario: given event should update entities
    Given I am admin
    When I send an event:
    """
    {
      "connector": "test-connector-che-3",
      "connector_name": "test-connector-name-che-3",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-che-3",
      "state": 2,
      "output": "test-output-che-3"
    }
    """
    When I save response createComponentTimestamp={{ now }}
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector": "test-connector-che-3",
      "connector_name": "test-connector-name-che-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-3",
      "resource": "test-resource-che-3",
      "state": 2,
      "output": "test-output-che-3"
    }
    """
    When I save response createResourceTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=che-3
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-component-che-3",
          "category": null,
          "component": "test-component-che-3",
          "depends": [
            "test-resource-che-3/test-component-che-3"
          ],
          "enable_history": [
            {{ .createComponentTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-connector-che-3/test-connector-name-che-3"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-component-che-3",
          "type": "component"
        },
        {
          "_id": "test-connector-che-3/test-connector-name-che-3",
          "category": null,
          "depends": [
            "test-component-che-3"
          ],
          "enable_history": [
            {{ .createComponentTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-resource-che-3/test-component-che-3"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-connector-name-che-3",
          "type": "connector"
        },
        {
          "_id": "test-resource-che-3/test-component-che-3",
          "category": null,
          "component": "test-component-che-3",
          "depends": [
            "test-connector-che-3/test-connector-name-che-3"
          ],
          "enable_history": [
            {{ .createResourceTimestamp }}
          ],
          "enabled": true,
          "impact": [
            "test-component-che-3"
          ],
          "impact_level": 1,
          "infos": {},
          "measurements": null,
          "name": "test-resource-che-3",
          "type": "resource"
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
