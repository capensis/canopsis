<template>
  <c-layout class="mt-3" gap="2" column>
    <span class="text-body-2">
      {{ t('common.column', 2) }}
    </span>

    <quick-alarm-actions-form
      :actions="actions"
      :massive="massive"
      name="actions"
      @update:actions="updateActions"
    />
  </c-layout>
</template>

<script>
import { watch } from 'vue';

import { useField, useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import QuickAlarmActionsForm from '@/components/common/actions-panel/quick-alarm-actions-form.vue';

export default {
  components: { QuickAlarmActionsForm },

  props: {
    massive: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const { t } = useI18n();
    const form = useFormContext();
    const actions = form.model('actions');

    const { validate } = useField(
      'actions',
      () => (actions.value.length > 0 ? true : t('widgetTemplate.errors.quickActionsRequired')),
      {
        initialValue: actions.value.length,
        validateOnValueUpdate: false,
      },
    );

    watch(
      () => actions.value.length,
      () => validate(),
      { flush: 'post' },
    );

    const updateActions = (value) => {
      actions.value = value;
    };

    return {
      t,
      actions,
      updateActions,
    };
  },
};
</script>
