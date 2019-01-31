import { createNamespacedHelpers } from 'vuex';
import { schema as normalizrSchema } from 'normalizr';

const { mapActions } = createNamespacedHelpers('entities');

const prepareUsingCountForArraySchema = (schema, value = [], oldValue = []) => {
  const usingCount = {};

  value.forEach((item) => {
    const id = item[schema.idAttribute];

    if (id) {
      usingCount[id] = (usingCount[id] || 0) + 1;
    }
  });

  oldValue.forEach((item) => {
    const id = item[schema.idAttribute];

    if (id) {
      usingCount[id] = (usingCount[id] || 0) - 1;
    }
  });

  return usingCount;
};

const prepareUsingCountForEntitySchema = (schema, value = {}, oldValue = {}) => {
  const usingCount = {};
  const newValueId = value[schema.idAttribute];
  const oldValueId = oldValue[schema.idAttribute];

  if (newValueId !== oldValueId) {
    if (newValueId) {
      usingCount[newValueId] = 1;
    }

    if (oldValueId) {
      usingCount[oldValueId] = -1;
    }
  }

  return usingCount;
};

export default (schema, entityFieldName) => {
  const isArray = schema instanceof normalizrSchema.Array || Array.isArray(schema);
  let entitySchema = schema;

  if (isArray) {
    entitySchema = schema[0] || schema.schema;
  }

  if (!entitySchema) {
    console.error('Incorrect entitySchema');

    return {};
  }

  return {
    watch: {
      [entityFieldName]: {
        immediate: true,
        handler(value, oldValue) {
          let usingCount = {};

          if (isArray) {
            usingCount = prepareUsingCountForArraySchema(schema, value, oldValue);
          } else {
            usingCount = prepareUsingCountForEntitySchema(schema, value, oldValue);
          }

          this.mergeEntitiesUsingCount({
            entity: entitySchema.key,
            usingCount,
          });
        },
      },
    },
    methods: {
      ...mapActions(['mergeEntitiesUsingCount']),
    },
  };
};
