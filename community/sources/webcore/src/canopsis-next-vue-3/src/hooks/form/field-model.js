import { computed, unref } from 'vue';
import { get } from 'lodash';

import { useFormContext } from 'vee-validate';

/**
 * Creates a writable model for a vee-validate form value
 * without registering a field.
 *
 * @param {string|import('vue').Ref<string>|import('vue').ComputedRef<string>} path
 * @return {import('vue').ComputedRef<*>}
 */
export const useFormModel = path => {
  const { values, setFieldValue } = useFormContext();

  return computed({
    get: () => get(values, unref(path)),
    set: value => {
      setFieldValue(unref(path), value);
    },
  });
};
