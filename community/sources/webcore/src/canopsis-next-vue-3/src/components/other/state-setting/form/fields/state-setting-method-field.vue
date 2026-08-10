<template>
  <c-layout gap="2" column>
    <c-label :label="t('stateSetting.computeMethod')" />
    <c-radio-group
      :model-value="modelValue"
      class="mt-0"
      hide-details
      mandatory
      @update:model-value="$emit('update:modelValue', $event)"
    >
      <v-radio
        v-for="method in STATE_SETTING_METHODS"
        :key="method"
        :value="method"
        color="primary"
      >
        <template #label>
          <span class="mr-2">
            {{ t(`stateSetting.methods.${method}.label`) }}
          </span>
          <c-help-icon :text="t(`stateSetting.methods.${method}.tooltip`)" icon="help" top />
        </template>
      </v-radio>
    </c-radio-group>
  </c-layout>
</template>

<script>
import { STATE_SETTING_METHODS } from '@/constants';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      default: STATE_SETTING_METHODS.inherited,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    return {
      t,
      STATE_SETTING_METHODS,
    };
  },
};
</script>
