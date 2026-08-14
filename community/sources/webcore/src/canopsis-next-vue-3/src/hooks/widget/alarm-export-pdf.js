import { inject } from 'vue';

import { exportAlarmToPdf as exportAlarmToPdfFile } from '@/helpers/file/pdf';
import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useAlarm } from '@/hooks/store/modules/alarm';

/**
 * Alarm PDF export with comments loaded from alarm details.
 *
 * @returns {Object}
 */
export const useAlarmExportPdf = () => {
  const { t } = useI18n();
  const popups = usePopups();
  const $system = inject('$system', {});
  const { fetchAlarmDetailsListWithoutStore } = useAlarm();

  const exportAlarmToPdf = async (alarm, template) => {
    try {
      const response = await fetchAlarmDetailsListWithoutStore({
        params: [
          {
            _id: alarm._id,
            steps: { type: 'comment' },
          },
        ],
      });

      const comments = response?.[0]?.data?.steps?.data ?? [];

      await exportAlarmToPdfFile(template, { ...alarm, comments }, $system.timezone);
    } catch (err) {
      console.error(err);
      popups.error({ text: t('errors.default') });
    }
  };

  return {
    exportAlarmToPdf,
  };
};
