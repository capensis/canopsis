import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

const { mapActions: serviceMapActions } = createNamespacedHelpers('service');
const { mapActions: dynamicInfoMapActions } = createNamespacedHelpers('dynamicInfo');

export const widgetColumnsInfos = {
  data() {
    return {
      alarmInfos: [],
      entityInfos: [],
      infosPending: false,
    };
  },
  mounted() {
    this.fetchInfos();
  },
  methods: {
    ...serviceMapActions({ fetchEntityInfosKeysWithoutStore: 'fetchInfosKeysWithoutStore' }),
    ...dynamicInfoMapActions({ fetchDynamicInfosKeysWithoutStore: 'fetchInfosKeysWithoutStore' }),

    async fetchInfos() {
      try {
        const params = { limit: MAX_LIMIT };

        this.infosPending = true;

        const [
          { data: alarmInfos },
          { data: entityInfos },
        ] = await Promise.all([
          this.fetchDynamicInfosKeysWithoutStore({ params }),
          this.fetchEntityInfosKeysWithoutStore({ params }),
        ]);

        this.alarmInfos = alarmInfos;
        this.entityInfos = entityInfos;
      } catch (err) {
        console.warn(err);
      } finally {
        this.infosPending = false;
      }
    },
  },
};
