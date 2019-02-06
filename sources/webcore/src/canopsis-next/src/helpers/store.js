import { get } from 'lodash';

import schemas from '@/store/schemas';

/**
 * Helper for entities preparation for delete
 *
 * @param type
 * @param data
 * @returns {{entitiesToMerge: {}, entitiesToDelete: {}}}
 */
export function prepareEntitiesToDelete({ type, data }) {
  const schema = schemas[type];
  const id = data[schema.idAttribute];

  const entitiesToMerge = {};
  const entitiesToDelete = {
    [schema.key]: {
      [id]: {},
    },
  };

  const prepareChild = (entity, childSchema) => {
    const parents = get(entity, '_embedded.parents', []);

    if (parents.length <= 1) {
      const result = prepareEntitiesToDelete({ type: childSchema.key, data: entity });

      Object.assign(entitiesToMerge, result.entitiesToMerge);
      Object.assign(entitiesToDelete, result.entitiesToDelete);
    } else {
      if (!entitiesToMerge[childSchema.key]) {
        entitiesToMerge[childSchema.key] = {};
      }

      entitiesToMerge[childSchema.key][entity[childSchema.idAttribute]] = {
        ...entity,
        _embedded: {
          ...entity._embedded,
          parents: parents.filter(v => v.id !== id || (v.id === id && v.type !== type)),
        },
      };
    }
  };

  Object.keys(schema.schema).forEach((key) => {
    if (Array.isArray(schema.schema[key])) {
      const childrenSchema = schema.schema[key][0];

      data[key].forEach(entity => prepareChild(entity, childrenSchema));
    } else {
      prepareChild(data[key], schema.schema[key]);
    }
  });

  return {
    entitiesToMerge,
    entitiesToDelete,
  };
}

export default {
  prepareEntitiesToDelete,
};
