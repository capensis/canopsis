# Associative Table in Canopsis

This document defines the specifications of the AssociativeTable feature of Canopsis 2.xx

## Goal

This feature aims to provide a standard way to store and read key/values pairs in the mongo database, identified with a specific key.

In fine, it can store configuration values for classes.

### Developer requirements

As a Canopsis developer, I want to store and read values with having to modify static configuration ".conf" files.

## Implementation

### Configuration management

The collection name where AssociativeTable are stored is given as a constant in AssociativeTableManager: **AssociativeTable.STORAGE_URI**

### Classes

class **AssociativeTable** :

 * Responsibility :
    * is a standard class to manipulate a key/value store


class **AssociativeTableManager** :

 * Responsibility :
     * instantiate AssociativeTable objects
     * read/store AssociativeTable objects in the database

## Usage

### Instantiate the manager

The manager need a logger and a pymongo collection objects. You can instantiate a new manager as follow:

```python
import logging
from canopsis.common.associative_table.manager import AssociativeTableManager
from pymongo import MongoClient

logger = logging.getLogger()
logger.setLevel(logging.DEBUG)
collection = MongoClient(**args)['my_collection_name']

my_manager = AssociativeTableManager(logger=logger, collection=collection)
```

Or, with a canopsis **Middleware** class:

```python
import logging
from canopsis.common.associative_table.manager import AssociativeTableManager
from canopsis.middleware.core import Middleware

logger = logging.getLogger()
logger.setLevel(logging.DEBUG)
storage = Middleware.get_middleware_by_uri(AssociativeTableManager.STORAGE_URI)

my_manager = AssociativeTableManager(logger=logger, collection=storage._backend)
```

### Read / Write
```python
assoc = my_manager.create('test')

assoc.set('john', 'cleese')

my_manager.save(assoc)

name = assoc.get('john')  # name == 'cleese'
```
