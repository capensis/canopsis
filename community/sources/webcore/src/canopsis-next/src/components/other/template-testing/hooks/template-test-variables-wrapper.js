import { get } from 'lodash';
import { computed, ref, unref, watch } from 'vue';

import { isOmitEqual } from '@/helpers/collection';
import { convertRuleToTemplateTestingTestValidateForm } from '@/helpers/entities/template-testing-validate/form';

import { useModelField } from '@/hooks/form/model-field';

/**
 * Hook for managing test variables fields with form integration
 *
 * @param {Object} [props={}] - Component props containing form data
 * @param {Object} props.form - Form object containing test variables data
 * @param {import('vue').Ref<string>|string} type - The type of template testing
 * @param {Function} emit - Vue emit function for component communication
 * @returns {Object} Hook return object
 * @returns {import('vue').Ref<Array>} returns.items - Reactive array of test variable field items with update handlers
 */
export const useTestVariablesFields = (props = {}, type, emit) => {
  const items = ref([]);

  const { updateField } = useModelField(props, emit);

  const isEmptyItems = computed(() => !items.value.length);

  watch(() => props.form, (newForm) => {
    const newItems = convertRuleToTemplateTestingTestValidateForm(newForm, unref(type));

    if (isOmitEqual(items.value, newItems, ['value', 'updateField'])) {
      return;
    }

    items.value = newItems.map((newItem) => {
      const key = newItem.formKey ?? newItem.key;

      return {
        ...newItem,
        value: get(newForm, key),
        updateField: value => updateField(key, value),
      };
    });
  }, { immediate: true });

  return {
    items,
    isEmptyItems,
  };
};
