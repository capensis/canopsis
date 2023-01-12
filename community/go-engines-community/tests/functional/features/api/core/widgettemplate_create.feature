Feature: Create a widget template
  I need to be able to create a widget template
  Only admin should be able to create a widget template

  Scenario: given create alarm request should return ok
    When I am admin
    When I do POST /api/v4/widget-templates:
    """json
    {
      "title": "test-widgettemplate-to-create-1-title",
      "type": "alarm",
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
      "type": "alarm",
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
      "type": "alarm",
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

  Scenario: given create entity request should return ok
    When I am admin
    When I do POST /api/v4/widget-templates:
    """json
    {
      "title": "test-widgettemplate-to-create-2-title",
      "type": "entity",
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
      "type": "entity",
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
      "type": "entity",
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
        "type": "Type is missing.",
        "columns": "Columns is missing."
      }
    }
    """
    When I do POST /api/v4/widget-templates:
    """json
    {
      "columns": []
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "columns": "Columns should not be blank."
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
        "type": "Type must be one of [alarm entity]."
      }
    }
    """

  Scenario: given create request with invalid columns should return bad request error
    When I am admin
    When I do POST /api/v4/widget-templates:
    """json
    {
      "title": "test-widgettemplate-to-create-3-title",
      "type": "alarm",
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
      "type": "alarm",
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
      "type": "entity",
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
