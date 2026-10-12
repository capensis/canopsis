<template>
  <v-tab
    :value="value"
    :disabled="disabled"
    :class="{ 'text-error': error }"
    class="v-tab--tooltip"
  >
    <c-tooltip :disabled="!disabled" location="top">
      <template #activator="{ props: activatorProps }">
        <span v-bind="activatorProps" class="test-variables-tab__wrapper">
          {{ t('templateTesting.testVariables') }}
        </span>
      </template>

      <span>{{ t('templateTesting.testVariablesDisabledTooltip') }}</span>
    </c-tooltip>
  </v-tab>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
    value: {
      type: [Number, String],
      required: true,
    },
    error: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const { t } = useI18n();

    return { t };
  },
};
</script>

<style lang="scss" scoped>
.test-variables-tab__wrapper {
  max-width: 360px;
}
</style>
