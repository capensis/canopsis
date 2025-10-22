import { get } from 'lodash';
import { computed, ref, unref, watch } from 'vue';

import { USER_PERMISSIONS } from '@/constants';

import { isOmitEqual } from '@/helpers/collection';
import { convertRuleToTemplateTestingTestValidateForm } from '@/helpers/entities/template-testing-validate/form';

import { useCanPermission } from '@/hooks/auth';
import { useModelField } from '@/hooks/form/model-field';

/**
 * Hook for managing test variables fields with form integration and access permission check
 *
 * @param {Object} [props={}] - Component props containing form data
 * @param {Object} props.form - Form object containing test variables data
 * @param {import('vue').Ref<string>|string} type - The type of template testing
 * @param {Function} emit - Vue emit function for component communication
 * @returns {Object} Hook return object
 * @returns {import('vue').Ref<Array>} returns.items - Reactive array of test variable field items with update handlers
 * @returns {import('vue').Ref<boolean>} returns.isEmptyItems - Reactive boolean indicating if the items array is empty
 * @returns {import('vue').Ref<boolean>} returns.hasAccess - Reactive boolean indicating if the user
 *                                                           has access to the template testing test variables
 */
export const useTestVariablesTabData = (props = {}, type, emit) => {
  const items = ref([]);

  const { updateField } = useModelField(props, emit);

  const { hasAccess } = useCanPermission(USER_PERMISSIONS.technical.templateTesting);

  const isEmptyItems = computed(() => !items.value.length);

  watch(() => props.form, (newForm) => {
    const newItems = convertRuleToTemplateTestingTestValidateForm(newForm, unref(type));

    if (isOmitEqual(items.value, newItems, ['value', 'updateField'])) {
      return;
    }

    items.value = newItems.map((newItem) => {
      const { key } = newItem;

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
    hasAccess,
  };
};
