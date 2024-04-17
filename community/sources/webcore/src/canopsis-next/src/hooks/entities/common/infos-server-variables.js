import { computed, unref } from 'vue';

import { PAYLOADS_INFO_VARIABLES } from '@/constants';

/**
 * Provides a reactive list of infos server variables
 *
 * @returns {Object}
 */
export const useInfosServerVariables = (infos = []) => {
  const subVariables = Object.values(PAYLOADS_INFO_VARIABLES).map(value => ({ value }));

  const variables = computed(() => {
    const fields = ['.%name%', ...unref(infos).map(({ value }) => `.${value}`)];

    return fields.map(value => ({
      value,
      variables: subVariables,
    }));
  });

  return {
    variables,
  };
};
