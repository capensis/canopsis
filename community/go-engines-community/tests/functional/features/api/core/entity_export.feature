Feature: Export entities
  I need to be able to export a entities

  Scenario: given export request should return entities
    When I am admin
    When I do POST /api/v4/entity-export:
    """json
    {
      "search": "test-entity-to-export",
      "fields": [
         {"name": "_id", "label": "ID"},
         {"name": "name", "label": "Name"},
         {"name": "type", "label": "Type"},
         {"name": "infos", "label": "Infos"}
      ]
    }
    """
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/entity-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/entity-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    ID,Name,Type,Infos
    test-entity-to-export-resource-1/test-entity-to-export-component,test-entity-to-export-resource-1,resource,"{""test-entity-to-export-resource-1-info-1"":{""name"":""test-entity-to-export-resource-1-info-1-name"",""description"":""test-entity-to-export-resource-1-info-1-description"",""value"":""test-entity-to-export-resource-1-info-1-value""}}"
    test-entity-to-export-resource-2/test-entity-to-export-component,test-entity-to-export-resource-2,resource,{}

    """

  Scenario: given export request should return empty response
    When I am admin
    When I do POST /api/v4/entity-export:
    """json
    {
      "search": "not found",
      "fields": [
         {"name": "_id", "label": "ID"},
         {"name": "name", "label": "Name"},
         {"name": "type", "label": "Type"},
         {"name": "infos", "label": "Infos"}
      ]
    }
    """
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/entity-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/entity-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    ID,Name,Type,Infos

    """


  Scenario: given not exit export task should return not found error
    When I am admin
    When I do GET /api/v4/entity-export/not-found
    Then the response code should be 404
    When I do GET /api/v4/entity-export/not-found/download
    Then the response code should be 404

  Scenario: given unauth request should not allow access
    When I do POST /api/v4/entity-export
    Then the response code should be 401
    When I do GET /api/v4/entity-export/not-found
    Then the response code should be 401
    When I do GET /api/v4/entity-export/not-found/download
    Then the response code should be 401

  Scenario: given request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/entity-export
    Then the response code should be 403
    When I do GET /api/v4/entity-export/not-found
    Then the response code should be 403
    When I do GET /api/v4/entity-export/not-found/download
    Then the response code should be 403
