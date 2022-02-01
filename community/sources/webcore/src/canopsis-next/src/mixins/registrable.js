import { createNamespacedHelpers } from 'vuex';
import { normalize } from 'normalizr';

import uid from '@/helpers/uid';

const { mapActions } = createNamespacedHelpers('entities');

const registerGetterMethodName = uid('registerGetter');
const unregisterGetterMethodName = uid('unregisterGetter');

/**
 * Create registrable mixin for components
 *
 * @param {Object} entitySchema - Normalizr schema of entity
 * @param {string} entityFieldName - Entity field name
 * @return {{methods: {}, beforeDestroy(): void, mounted(): void}}
 */
export const registrableMixin = (entitySchema, entityFieldName) => {
  const dependenciesPreparer = (entities) => {
    const { entities: normalizedEntities } = normalize(entities, entitySchema);

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
