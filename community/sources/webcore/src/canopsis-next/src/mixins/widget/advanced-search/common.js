import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

const { mapActions: mapServiceActions } = createNamespacedHelpers('service');

export const widgetAdvancedSearchCommonMixin = {
  data() {
    return {
      entityInfosKeys: [],
      entityInfosKeysPending: false,
    };
  },
  mounted() {
    this.fetchInfosKeysList();
  },
  methods: {
    ...mapServiceActions({ fetchEntityInfosKeysWithoutStore: 'fetchInfosKeysWithoutStore' }),

    getDeepInfosItems(valuePrefix = '', textPrefix = '') {
      return [
        {
          value: `${valuePrefix}.name`,
          text: this.$t('common.name'),
          selectorText: `${textPrefix}.${this.$t('common.name')}`,
        }, {
          value: `${valuePrefix}.value`,
          text: this.$t('common.value'),
          selectorText: `${textPrefix}.${this.$t('common.value')}`,
        },
      ];
    },

    async fetchInfosKeysList() {
      try {
        this.entityInfosKeysPending = true;

        const { data: infos } = await this.fetchEntityInfosKeysWithoutStore({
          params: { limit: MAX_LIMIT },
        });

        this.entityInfosKeys = infos;
      } catch (err) {
        console.error(err);
      } finally {
        this.entityInfosKeysPending = false;
      }
    },
  },
};
