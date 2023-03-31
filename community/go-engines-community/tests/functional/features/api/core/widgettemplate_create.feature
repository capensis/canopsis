Feature: Create a widget template
  I need to be able to create a widget template
  Only admin should be able to create a widget template

  Scenario: given create alarm columns request should return ok
    When I am admin
    When I do POST /api/v4/widget-templates:
    """json
    {
      "title": "test-widgettemplate-to-create-1-title",
      "type": "alarm_columns",
      "columns": [
        {
          "value": "v.resource"
        },
        {
          "value": "v.component"
        },
        {
          "value": "extra_details"
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-widgettemplate-to-create-1-title",
      "type": "alarm_columns",
      "columns": [
        {
          "value": "v.resource"
        },
        {
          "value": "v.component"
        },
        {
          "value": "extra_details"
        }
      ]
    }
    """
    When I do GET /api/v4/widget-templates/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-widgettemplate-to-create-1-title",
      "type": "alarm_columns",
      "columns": [
        {
          "value": "v.resource"
        },
        {
          "value": "v.component"
        },
        {
          "value": "extra_details"
        }
      ]
    }
    """

  Scenario: given create entity columns request should return ok
    When I am admin
    When I do POST /api/v4/widget-templates:
    """json
    {
      "title": "test-widgettemplate-to-create-2-title",
      "type": "entity_columns",
      "columns": [
        {
          "value": "_id"
        },
        {
          "value": "type"
        },
        {
          "value": "infos.test.value"
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-widgettemplate-to-create-2-title",
      "type": "entity_columns",
      "columns": [
        {
          "value": "_id"
        },
        {
          "value": "type"
        },
        {
          "value": "infos.test.value"
        }
      ]
    }
    """
    When I do GET /api/v4/widget-templates/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-widgettemplate-to-create-2-title",
      "type": "entity_columns",
      "columns": [
        {
          "value": "_id"
        },
        {
          "value": "type"
        },
        {
          "value": "infos.test.value"
        }
      ]
    }
    """

  Scenario: given create alarm more infos request should return ok
    When I am admin
    When I do POST /api/v4/widget-templates:
    """json
    {
      "title": "test-widgettemplate-to-create-3-title",
      "type": "alarm_more_infos",
      "content": "{{ `{{ alarm.v.display_name }}` }}"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-widgettemplate-to-create-3-title",
      "type": "alarm_more_infos",
      "content": "{{ `{{ alarm.v.display_name }}` }}"
    }
    """
    When I do GET /api/v4/widget-templates/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-widgettemplate-to-create-3-title",
      "type": "alarm_more_infos",
      "content": "{{ `{{ alarm.v.display_name }}` }}"
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/widget-templates
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/widget-templates
    Then the response code should be 403

  Scenario: given create request with missing fields should return bad request error
    When I am admin
    When I do POST /api/v4/widget-templates:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "title": "Title is missing.",
        "type": "Type is missing."
      }
    }
    """
    When I do POST /api/v4/widget-templates:
    """json
    {
      "type": "alarm_columns"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "columns": "Columns is missing."
      }
    }
    """
    When I do POST /api/v4/widget-templates:
    """json
    {
      "type": "alarm_columns",
      "columns": []
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "columns": "Columns is missing."
      }
    }
    """
    When I do POST /api/v4/widget-templates:
    """json
    {
      "type": "alarm_more_infos"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "content": "Content is missing."
      }
    }
    """
    When I do POST /api/v4/widget-templates:
    """json
    {
      "type": "alarm_more_infos",
      "content": ""
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "content": "Content is missing."
      }
    }
    """
    When I do POST /api/v4/widget-templates:
    """json
    {
      "type": "unknown"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "type": "Type must be one of [alarm_columns entity_columns alarm_more_infos weather_item weather_modal weather_entity]."
      }
    }
    """

  Scenario: given create request with invalid columns should return bad request error
    When I am admin
    When I do POST /api/v4/widget-templates:
    """json
    {
      "title": "test-widgettemplate-to-create-3-title",
      "type": "alarm_columns",
      "columns": [
        {}
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "columns.0.value": "Value is missing."
      }
    }
    """
    When I do POST /api/v4/widget-templates:
    """json
    {
      "title": "test-widgettemplate-to-create-1-title",
      "type": "alarm_columns",
      "columns": [
        {
          "value": "unknown"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "columns.0.value": "Value is invalid."
      }
    }
    """
    When I do POST /api/v4/widget-templates:
    """json
    {
      "type": "entity_columns",
      "columns": [
        {
          "value": "unknown"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "columns.0.value": "Value is invalid."
      }
    }
    """
