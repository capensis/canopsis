import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

const { mapActions: mapAlarmActions } = createNamespacedHelpers('alarm');

export const patternCountAlarmsModalMixin = {
  methods: {
    ...mapAlarmActions({ fetchAlarmsListWithoutStore: 'fetchListWithoutStore' }),

    showAlarmsModalByPatterns(patterns) {
      const widget = generatePreparedDefaultAlarmListWidget();

      this.$modals.show({
        name: MODALS.alarmsList,
        config: {
          widget,
          title: this.$t('pattern.patternAlarms'),
          fetchList: params => this.fetchAlarmsListWithoutStore({
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
