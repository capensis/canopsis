import { computed, ref, unref } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { useCopyVars } from '@/hooks/store/modules/copy-vars';
import { usePendingHandler } from '@/hooks/query/pending';

import { useVarsPrepare } from './shared';

/**
 * Hook for managing copy variables list based on test type
 *
 * @param {Object} options - Configuration options
 * @param {import('vue').Ref<string>|string} [options.type] - The type of template testing
 *                                                            (eventFilter, linkRule, scenario, etc.)
 * @param {import('vue').Ref<Object>|Object} [options.form] - The form object
 *
 * @returns {Object} Hook return object
 * @returns {import('vue').Ref<Array>} returns.templateVars - Reactive array of template variables
 * @returns {import('vue').Ref<boolean>} returns.pending - Loading state indicator
 * @returns {Function} returns.fetchList - Function to fetch template variables list
 */
export const useCopyVarsList = ({ type } = {}) => {
  const vars = ref({});

  const {
    fetchEventFiltersCopyVarsWithoutStore,
    fetchDynamicInfosCopyVarsWithoutStore,
  } = useCopyVars();

  const fetchHandler = computed(() => ({
    [TEMPLATE_TESTING_TEST_TYPES.eventFilter]: fetchEventFiltersCopyVarsWithoutStore,
    [TEMPLATE_TESTING_TEST_TYPES.dynamicInfo]: fetchDynamicInfosCopyVarsWithoutStore,
  }[unref(type)]));

  const { prepare: prepareVars } = useVarsPrepare('templateTesting.templateVars');

  const {
    pending,
    handler: fetchList,
  } = usePendingHandler(async () => {
    if (!fetchHandler.value) {
      return;
    }

    const data = await fetchHandler.value();

    vars.value = prepareVars(data);
  });

  return {
    vars,
    pending,
    fetchList,
  };
};
