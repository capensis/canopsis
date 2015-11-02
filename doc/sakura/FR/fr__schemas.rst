.. _FR__Schemas:

=======
Schemas
=======

This document presents schema specification.

.. contents::
   :depth: 2

References
==========


Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Gwenael Pluchon", "2015/11/02", "0.1", "", ""

Contents
========

Meta-Schema
-----------

{
  "$schema": "meta-schema",
  "type": "object",
  "required": [ "properties" ],
  "properties": {
    "properties": {
      "type": "object"
    },

    "categories": {
      "type": "array",
      "items": {
        "type": "object"
        "properties": {
           "title": { "type": "string" }
           "keys": {
              "type": "array",
              "items": { "type": "string" }
            }
        }
      }
    },

    "metadata": {
      "type": "object"
      "properties": {
        "description": { "type": "string" },
        "icon": { "type": "string" }
      }
    }
  }
}
