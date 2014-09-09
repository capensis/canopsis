Composite Storage
=================

.. module:: canopsis.storage.composite

.. class:: CompositeStorage(canopsis.storage.Storage)

   Storage dedicated to manage composite data.

   .. data:: __datatype__ = 'composite'

      storage data type name

   .. method:: get(_ids=None, data_type=None, limit=0, skip=0, sort=None)

      Get a list of data identified among data_ids or a type

      :param data_ids: data ids to get
      :type data_id: list of str

      :param data_type: data_id type to get if not None
      :type data_type: str

      :param limit: max number of data to get
      :type limit: int

      :param skip: starting index of research if multi data to get
      :type skip: int

      :param sort: couples of field (name, value) to sort with ASC/DESC
         Storage fields
      :type sort: dict

      :return: a list of couples of field (name, value) or None respectivelly
         if such data exist or not
      :rtype: list of dict of field (name, value)

   .. method:: put(_id, data, data_type=None)

      Put a data related to an id

      :param _id: data id
      :type _id: str

      :param data_type: data type to update
      :type data_type: str

      :param data: data to update
      :type data: dict

   .. method:: remove(_ids=None, data_type=None):

      Remove data from ids or type

      :param _ids: list of data id
      :type _ids: list

      :param data_type: data type to remove if not None
      :type data_type: str
