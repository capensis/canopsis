<template>
  <c-select
    :model-value="value"
    :items="items"
    :label="label || t('common.filterByType')"
    :name="name"
    :disabled="disabled"
    :required="required"
    :clearable="clearable"
    :hide-details="hideDetails"
    class="template-testing-data-type-field"
    item-title="text"
    item-value="value"
    @update:model-value="$emit('update:value', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import { TEMPLATE_TESTING_DATA_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    value: {
      type: Number,
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'type',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:value'],
  setup() {
    const { t } = useI18n();

    const items = computed(() => Object.values(TEMPLATE_TESTING_DATA_TYPES).map(value => ({
      value,
      text: t(`templateTesting.dataTypes.${value}`),
    })));

    return {
      t,
      items,
    };
  },
};
</script>

<style lang="scss" scoped>
.template-testing-data-type-field {
  min-width: 200px;
}
</style>
