Feature: Upload a MIB file
  I need to be able to upload a MIB file
  Only admin should be able to upload a MIB file

  Scenario: given upload request and no auth user should not allow access
    When I do POST /api/v4/cat/snmpmibs
    Then the response code should be 401

  Scenario: given upload request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/snmpmibs
    Then the response code should be 403

  Scenario: given upload request with wrong MIB should return Bad Request
    When I am admin
    When I set header Content-Type=application/x-www-form-urlencoded
    When I do POST /api/v4/cat/snmpmibs:
    """
    filecontent=%5B%7B%22filename%22%3A%20%22concatenatedMibFiles%22%2C%22data%22%3A%20%22WRONGMIBDARLING%22%7D%5D
    """
    Then the response code should be 400

  Scenario: given upload request with valid MIB should return OK result
    When I am admin
    When I set header Content-Type=application/x-www-form-urlencoded
    When I do POST /api/v4/cat/snmpmibs:
    """
    filecontent=%5B%7B%22filename%22%3A%20%22MibFile%22%2C%22data%22%3A%20%22WRONGMIBDARLING%22%7D%5D
    """
    Then the response code should be 200
