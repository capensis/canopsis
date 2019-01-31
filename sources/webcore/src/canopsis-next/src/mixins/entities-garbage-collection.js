import { createNamespacedHelpers } from 'vuex';
import { normalize } from 'normalizr';

import uid from '@/helpers/uid';

const { mapActions } = createNamespacedHelpers('entities');

const registerGetterMethodName = uid('registerGetter');
const unregisterGetterMethodName = uid('unregisterGetter');

export default (schema, entityFieldName) => {
  const dependenciesPreparer = (entities) => {
    const { entities: normalizedEntities } = normalize(entities, schema);

    return Object.keys(normalizedEntities).map(type => ({
      type,
      ids: Object.keys(normalizedEntities[type]),
    }));
  };

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
