import { createNamespacedHelpers } from 'vuex';
import { schema as normalizrSchema } from 'normalizr';

import uid from '@/helpers/uid';

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

const markMethodName = uid('mark');
const sweepMethodName = uid('sweep');

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
            usingCount = prepareUsingCountForArraySchema(entitySchema, value, oldValue);
          } else {
            usingCount = prepareUsingCountForEntitySchema(entitySchema, value, oldValue);
          }

          this[markMethodName]({
            type: entitySchema.key,
            usingCount,
          });
        },
      },
    },
    async beforeDestroy() {
      let usingCount;

      if (isArray) {
        usingCount = prepareUsingCountForArraySchema(entitySchema, [], this[entityFieldName]);
      } else {
        usingCount = prepareUsingCountForEntitySchema(entitySchema, {}, this[entityFieldName]);
      }

      await this[markMethodName]({
        type: entitySchema.key,
        usingCount,
      });

      this[sweepMethodName]({ type: entitySchema.key });
    },
    methods: {
      ...mapActions({
        [markMethodName]: 'mark',
        [sweepMethodName]: 'sweep',
      }),
    },
  };
};
