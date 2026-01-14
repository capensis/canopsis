import { PATTERNS_FIELDS, MODALS } from '@/constants';

import { formGroupsToPatternRulesQuery } from '@/helpers/entities/pattern/form';
import { generatePreparedDefaultContextWidget } from '@/helpers/entities/widget/form';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook for managing pattern count entities modal functionality.
 *
 * @param {Object} props - Component props containing pattern configuration flags.
 * @returns {Object} Object containing showPatternEntitiesModal function.
 * @property {Function} showPatternEntitiesModal - Function to show the entities modal with pattern filters.
 */
export const usePatternCountEntitiesModal = (props) => {
  const modals = useModals();
  const { t } = useI18n();
  const { useActions } = useStoreModuleHooks('entity');
  const { fetchContextEntitiesWithoutStore } = useActions({ fetchContextEntitiesWithoutStore: 'fetchListWithoutStore' });

  /**
   * Shows a modal with entities filtered by the specified pattern names.
   *
   * Only includes patterns that are enabled based on props configuration.
   *
   * @param {Array<string>} [patternNames=[
   *   PATTERNS_FIELDS.entity,
   *   PATTERNS_FIELDS.pbehavior,
   *   PATTERNS_FIELDS.event,
   *   PATTERNS_FIELDS.totalEntity,
   *   PATTERNS_FIELDS.serviceWeather
   * ]]
   */
  const showPatternEntitiesModal = (
    patternNames = [
      PATTERNS_FIELDS.entity,
      PATTERNS_FIELDS.pbehavior,
      PATTERNS_FIELDS.event,
      PATTERNS_FIELDS.totalEntity,
      PATTERNS_FIELDS.serviceWeather,
    ],
  ) => {
    const withPatternsByNamesMap = {
      [PATTERNS_FIELDS.entity]: true,
      [PATTERNS_FIELDS.pbehavior]: props.withPbehavior,
      [PATTERNS_FIELDS.event]: props.eventName,
      [PATTERNS_FIELDS.totalEntity]: props.totalEntityName,
      [PATTERNS_FIELDS.serviceWeather]: props.serviceWeatherName,
    };

    const patterns = patternNames.reduce((acc, patternName) => {
      if (!withPatternsByNamesMap[patternName]) {
        return acc;
      }

      acc[patternName] = formGroupsToPatternRulesQuery(props.value[patternName]?.groups);

      return acc;
    }, {});

    const widget = generatePreparedDefaultContextWidget();

    modals.show({
      name: MODALS.entitiesList,
      config: {
        widget,
        title: t('pattern.patternEntities'),
        fetchList: params => fetchContextEntitiesWithoutStore({
          params: {
            ...params,
            ...patterns,
          },
        }),
      },
    });
  };

  return { showPatternEntitiesModal };
};
