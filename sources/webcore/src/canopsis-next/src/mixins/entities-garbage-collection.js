import { createNamespacedHelpers } from 'vuex';
import { schema as normalizrSchema } from 'normalizr';

import uid from '@/helpers/uid';

const { mapActions } = createNamespacedHelpers('entities');

const registerGetterMethodName = uid('registerGetter');
const unregisterGetterMethodName = uid('unregisterGetter');

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

  let dependenciesPreparer;

  if (isArray) {
    dependenciesPreparer = entities => entities.map(entity => ({ type: entitySchema.key, id: entity._id }));
  } else {
    dependenciesPreparer = entity => ({ type: entitySchema.key, id: entity._id });
  }

  return {
    mounted() {
      this[registerGetterMethodName]({
        getDependencies: () => dependenciesPreparer(this[entityFieldName]),
        instance: this,
      });
    },
    beforeDestroy() {
      this[unregisterGetterMethodName](this);
    },
    methods: {
      ...mapActions({
        [registerGetterMethodName]: 'registerGetter',
        [unregisterGetterMethodName]: 'unregisterGetter',
      }),
    },
  };
};
