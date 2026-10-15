import { createNamespacedHelpers } from 'vuex';

import { exportAlarmToPdf } from '@/helpers/file/pdf';

const { mapActions } = createNamespacedHelpers('alarm/details');

export const widgetActionsPanelAlarmExportPdfMixin = {
  inject: ['$system'],
  methods: {
    ...mapActions({
      fetchAlarmDetailsWithoutStore: 'fetchListWithoutStore',
    }),

    async exportAlarmToPdf(alarm, template) {
      try {
        const response = await this.fetchAlarmDetailsWithoutStore({
          params: [{
            _id: alarm._id,
            steps: { type: 'comment' },
          }],
        });

        const comments = response?.[0]?.data?.steps?.data ?? [];
        const alarmWithComments = { ...alarm, comments };

        await exportAlarmToPdf(template, alarmWithComments, this.$system.timezone);
      } catch (err) {
        console.error(err);
        this.$popups.error({ text: this.$t('errors.default') });
      }
    },
  },
};
