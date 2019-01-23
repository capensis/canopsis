import { createNamespacedHelpers } from 'vuex';

import schemas from '@/store/schemas';

const { mapActions } = createNamespacedHelpers('entities');

export default (entity, computedTitle) => {
  const schema = schemas[entity];

  return {
    watch: {
      [computedTitle](value, oldValue) {
        const accumulator = {};

        value.forEach((item) => {
          const id = item[schema.idAttribute];

          if (id) {
            if (!accumulator[id]) {
              accumulator[id] = 0;
            }

            accumulator[id] += 1;
          }
        });

        oldValue.forEach((item) => {
          const id = item[schema.idAttribute];

          if (id) {
            if (!accumulator[id]) {
              accumulator[id] = 0;
            }

            accumulator[id] -= 1;
          }
        });

        this.mergeEntitiesUsingCount({ entity, usingCount: accumulator });
      },
    },
    methods: {
      ...mapActions(['mergeEntitiesUsingCount']),
    },
  };
};
