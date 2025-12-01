import { isNull } from 'lodash';
import { computed, ref, unref } from 'vue';

import { MODALS, PATTERN_OPERATORS } from '@/constants';

import { useModals } from '@/hooks/modals';
import { usePendingHandler } from '@/hooks/query/pending';

export const usePatternOptimization = (value) => {
  const modals = useModals();

  const status = ref(null);
  const suggestions = ref([]);

  const isOptimizing = computed(() => !isNull(status.value));

  const hasRegexpPatterns = computed(() => (
    Object.values(unref(value)).some(({ groups = [] }) => (
      groups.some(group => (
        group.rules.some(rule => rule.operator === PATTERN_OPERATORS.regexp)
      ))
    ))
  ));

  // TODO: replace it to usePolling hook for status checking
  const { pending, handler: tryOptimization } = usePendingHandler(async () => {
    const response = {
      status: 2,
      suggestions: [
        {
          entity_pattern: [
            [
              {
                field: 'infos.CSU_ID',
                field_type: 'string',
                cond: {
                  type: 'eq',
                  value: 'DO0011TVD0',
                },
              },
            ],
          ],
          found_entities: 218,
          difference: 0,
        },
        {
          entity_pattern: [
            [
              {
                field: 'component',
                cond: {
                  type: 'eq',
                  value: 'DO0011TVD0',
                },
              },
            ],
          ],
          found_entities: 1,
          difference: 217,
        },
        {
          entity_pattern: [
            [
              {
                field: 'name',
                cond: {
                  type: 'eq',
                  value: 'DO0011TVD0',
                },
              },
            ],
          ],
          found_entities: 1,
          difference: 217,
        },
        {
          entity_pattern: [
            [
              {
                field: 'infos.ticket_on_asset',
                field_type: 'string',
                cond: {
                  type: 'eq',
                  value: 'DO0011TVD0',
                },
              },
            ],
          ],
          found_entities: 1,
          difference: 217,
        },
      ],
    };

    status.value = response.status;
    suggestions.value = response.suggestions;
  });

  /**
     * Opens entities comparison modal for pattern optimization
     */
  const showEntitiesComparisonModal = () => modals.show({
    name: MODALS.entitiesComparison,
  });

  /**
   * Cancels pattern optimization process
   */
  const cancelOptimization = () => {
    // TODO: implement cancel optimization logic
  };

  return {
    isOptimizing,
    pending,
    suggestions,
    failed: false,
    hasRegexpPatterns,
    showEntitiesComparisonModal,
    tryOptimization,
    cancelOptimization,
  };
};
