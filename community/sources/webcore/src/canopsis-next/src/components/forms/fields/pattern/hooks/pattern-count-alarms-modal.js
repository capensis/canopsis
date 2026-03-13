import { MODALS } from '@/constants';

import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useAlarm } from '@/hooks/store/modules/alarm';

/**
 * Hook for showing alarms modal filtered by patterns
 *
 * @returns {Object} An object containing the showAlarmsModalByPatterns function
 */
export const usePatternCountAlarmsModal = () => {
  const modals = useModals();
  const { t } = useI18n();
  const { fetchAlarmsListWithoutStore } = useAlarm();

  /**
   * Shows alarms modal filtered by patterns
   *
   * @param {Object} patterns - Patterns object containing alarm_pattern, entity_pattern,
   *                            pbehavior_pattern and other pattern filters
   */
  const showAlarmsModalByPatterns = (patterns) => {
    const widget = generatePreparedDefaultAlarmListWidget();

    modals.show({
      name: MODALS.alarmsList,
      config: {
        widget,
        title: t('pattern.patternAlarms'),
        fetchList: params => fetchAlarmsListWithoutStore({
          params: {
            ...params,
            ...patterns,
          },
        }),
      },
    });
  };

  return {
    showAlarmsModalByPatterns,
  };
};
