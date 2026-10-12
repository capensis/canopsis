import { PATTERNS_FIELDS, MODALS } from '@/constants';

import { formGroupsToPatternRulesQuery } from '@/helpers/entities/pattern/form';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

import { useStoreModuleHooks } from '@/hooks/store';
import { useModals } from '@/hooks/modals';
import { useI18n } from '@/hooks/i18n';

/**
 * Hook for managing pattern count alarms modal functionality.
 *
 * @param {import('vue').Ref<Object>} value - Reactive reference to the pattern value object.
 * @returns {Object} Object containing showPatternAlarmsModal function.
 * @property {Function} showPatternAlarmsModal - Function to show the alarms modal with pattern filters.
 */
export const usePatternCountAlarmsModal = (props) => {
  const modals = useModals();
  const { t } = useI18n();
  const { useActions } = useStoreModuleHooks('alarm');
  const { fetchAlarmsListWithoutStore } = useActions({ fetchAlarmsListWithoutStore: 'fetchListWithoutStore' });

  /**
   * Shows a modal with alarms filtered by the specified pattern names.
   *
   * @param {Array<string>} [patternNames=[PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity, PATTERNS_FIELDS.pbehavior]]
   */
  const showPatternAlarmsModal = (
    patternNames = [
      PATTERNS_FIELDS.alarm,
      PATTERNS_FIELDS.entity,
      PATTERNS_FIELDS.pbehavior,
    ],
  ) => {
    const patterns = patternNames.reduce((acc, patternName) => {
      acc[patternName] = formGroupsToPatternRulesQuery(props.value[patternName]?.groups);

      return acc;
    }, {});

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

  return { showPatternAlarmsModal };
};
