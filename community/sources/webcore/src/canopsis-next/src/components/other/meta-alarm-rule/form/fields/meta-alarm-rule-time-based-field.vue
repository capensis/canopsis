<template>
  <c-form-block-row :label="t('metaAlarmRule.timeInterval')">
    <c-layout gap="2" align="center">
      <c-duration-field
        :duration="modelValue"
        :label="t('metaAlarmRule.timeInterval')"
        name="time_interval"
        required
        @update:duration="updateModel"
      />
      <c-help-icon
        :text="t('metaAlarmRule.timeIntervalHelpText')"
        icon="help"
        max-width="300"
        left
      />
    </c-layout>
  </c-form-block-row>
</template>

<script>
import { durationToForm } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: Object,
      default: () => durationToForm(),
    },
  },

  emits: ['update:modelValue'],

  setup(_, { emit }) {
    const { t } = useI18n();
    const updateModel = value => emit('update:modelValue', value);

    return {
      t,
      updateModel,
    };
  },
};
</script>
