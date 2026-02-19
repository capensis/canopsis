import { MODALS } from '@/constants';

import { generatePreparedDefaultContextWidget } from '@/helpers/entities/widget/form';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useEntity } from '@/hooks/store/modules/entity';

/**
 * Hook for showing entities modal filtered by patterns
 *
 * @returns {Object} An object containing the showEntitiesModalByPatterns function
 */
export const usePatternCountEntitiesModal = () => {
  const modals = useModals();
  const { t } = useI18n();
  const { fetchContextEntitiesListWithoutStore } = useEntity();

  /**
   * Shows entities modal filtered by patterns
   *
   * @param {Object} patterns - Patterns object containing entity_pattern and other pattern filters
   */
  const showEntitiesModalByPatterns = (patterns) => {
    const widget = generatePreparedDefaultContextWidget();

    modals.show({
      name: MODALS.entitiesList,
      config: {
        widget,
        title: t('pattern.patternEntities'),
        fetchList: params => fetchContextEntitiesListWithoutStore({
          params: {
            ...params,
            ...patterns,
          },
        }),
      },
    });
  };

  return {
    showEntitiesModalByPatterns,
  };
};
