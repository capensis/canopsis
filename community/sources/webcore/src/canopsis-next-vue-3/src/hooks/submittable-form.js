import { computed, unref } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { usePendingHandler } from '@/hooks/query/pending';

/**
 * Creates a submittable form handler with validation and error handling.
 *
 * @param {Object} options - Configuration options for the submittable form.
 * @param {Function} options.method - The submission method to be called.
 * @param {Object} [options.item] - Current item.
 * @param {Function} [options.isNewCheck = value => !value?._id] - The function to check if the item is new.
 * @returns {Object} An object containing methods and properties to manage the form submission.
 */
export const useSubmittableForm = ({ item, method, isNewCheck = value => !value?._id } = {}) => {
  const { t } = useI18n();

  const { pending: submitting, handler: submit } = usePendingHandler(method);

  const isNew = computed(() => isNewCheck?.(unref(item)));

  const submitLabel = computed(() => (isNew.value ? t('common.create') : t('common.save')));

  return {
    submitting,
    submit,
    isNew,
    submitLabel,
  };
};
