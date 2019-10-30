import { get, merge } from 'lodash';

import schemas from '@/store/schemas';
import { SCHEMA_EMBEDDED_KEY } from '@/config';

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

  let entitiesToMerge = {};
  let entitiesToDelete = {
    [schema.key]: {
      [id]: {},
    },
  };

  const prepareChild = (entity, childSchema) => {
    const parents = get(entity, `${SCHEMA_EMBEDDED_KEY}.parents`, []);

    if (parents.length <= 1) {
      const result = prepareEntitiesToDelete({ type: childSchema.key, data: entity });

      entitiesToMerge = merge(entitiesToMerge, result.entitiesToMerge);
      entitiesToDelete = merge(entitiesToDelete, result.entitiesToDelete);
    } else {
      if (!entitiesToMerge[childSchema.key]) {
        entitiesToMerge[childSchema.key] = {};
      }

      entitiesToMerge[childSchema.key][entity[childSchema.idAttribute]] = {
        ...entity,

        [SCHEMA_EMBEDDED_KEY]: {
          ...entity[SCHEMA_EMBEDDED_KEY],

          parents: parents.filter(v => v.id !== id || (v.id === id && v.type !== type)),
        },
      };
    }
  };

  Object.keys(schema.schema).forEach((key) => {
    if (data[key]) {
      if (Array.isArray(schema.schema[key])) {
        const childrenSchema = schema.schema[key][0];

        data[key].forEach(entity => prepareChild(entity, childrenSchema));
      } else {
        prepareChild(data[key], schema.schema[key]);
      }
    }
  });

  return {
    entitiesToMerge,
    entitiesToDelete,
  };
}

/**
 * Create deep clone of schema instance with `SCHEMA_EMBEDDED_KEY=true` for returning `SCHEMA_EMBEDDED_KEY` key
 *
 * @param schema
 * @return {any}
 */
export function cloneSchemaWithEmbedded(schema) {
  const newSchema = Object.assign(Object.create(Object.getPrototypeOf(schema)), schema);

  newSchema.schema = Object.keys(newSchema.schema).reduce((acc, key) => {
    if (Array.isArray(newSchema.schema[key])) {
      acc[key] = [cloneSchemaWithEmbedded(newSchema.schema[key][0])];
    } else {
      acc[key] = cloneSchemaWithEmbedded(newSchema.schema[key]);
    }

    return acc;
  }, {});

  newSchema[SCHEMA_EMBEDDED_KEY] = true;

  return newSchema;
}
