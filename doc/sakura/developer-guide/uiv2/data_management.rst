Data management
===============

1. Adapter
----------

We redefined several methods to recover data in application store because type
in Mongo database is 'view' and not 'userview'.

    * find :
    * findAll :


2. Serializer
-------------
We have defined and redefined any methods :

    * normalizeEmbeddedRelationships :
        - parameters :
            + type: type of main model
            + hash collected by ajax query : { data : array( records ),
                                               success : true }
        - returns normalized hash : { data : array(normalized records),
                                      embedded : array( normalized embedded records ),
                                      success : true }
        - manages hasMany/belongsTo relationships and polymorphic option


    * Redefined methods :
        - normalizeId :
          The redefined method takes as parameter the type of the hash and uses
          it to generate id.

        - normalize :
          The method has been redefined to take into account the modification on
          normalizeId

        - normalizePayload :
          By default, this method returns data in entry without traitment.
          Redefined method calls normalizeId and normalizeEmbeddedRelationships
          to format payload.

        - extractFindAll :
          Redefined method calls normalizePayload to normalize data in entry.
          then it pushes in store embedded records. Finally, method returns
          normalized main records to the store.