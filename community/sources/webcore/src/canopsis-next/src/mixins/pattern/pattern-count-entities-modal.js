import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import { generatePreparedDefaultContextWidget } from '@/helpers/entities/widget/form';

const { mapActions: mapEntityActions } = createNamespacedHelpers('entity');

export const patternCountEntitiesModalMixin = {
  methods: {
    ...mapEntityActions({ fetchContextEntitiesWithoutStore: 'fetchListWithoutStore' }),

    showEntitiesModalByPatterns(patterns) {
      const widget = generatePreparedDefaultContextWidget();

      this.$modals.show({
        name: MODALS.entitiesList,
        config: {
          widget,
          title: this.$t('pattern.patternEntities'),
          fetchList: params => this.fetchContextEntitiesWithoutStore({
            params: {
              ...params,
              ...patterns,
            },
          }),
        },
      });
    },
  },
};
