<template>
  <c-form-block-row :label="t('metaAlarmRule.childInactiveDelay')">
    <c-layout gap="3" column>
      <c-layout gap="2">
        <c-enabled-field
          :model-value="enabled"
          :label="t('metaAlarmRule.childInactiveDelay')"
          hide-details
          @update:model-value="toggle"
        />
        <c-help-icon
          :text="t('metaAlarmRule.childInactiveDelayHelpText')"
          icon="help"
          max-width="300"
          top
        />
      </c-layout>

      <c-duration-field
        v-if="modelValue"
        :duration="modelValue"
        name="child_inactive_delay"
        clearable
        @update:duration="updateModel"
      />
    </c-layout>
  </c-form-block-row>
</template>

<script>
import { computed } from 'vue';

import { durationToForm } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: Object,
      default: undefined,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const enabled = computed(() => !!props.modelValue);

    const updateModel = value => emit('update:modelValue', value);
    const toggle = value => updateModel(value ? durationToForm() : undefined);

    return {
      t,
      enabled,
      updateModel,
      toggle,
    };
  },
};
</script>
