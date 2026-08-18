<template>
  <c-select
    :model-value="value"
    :items="items"
    :label="label || t('templateTesting.filterByRuleType')"
    :name="name"
    :clearable="clearable"
    :required="required"
    :disabled="disabled"
    :hide-details="hideDetails"
    class="template-testing-test-type-field"
    item-title="text"
    item-value="value"
    @update:model-value="$emit('update:value', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import {
  TEMPLATE_TESTING_TEST_TYPES,
  TEMPLATE_TESTING_TESTS_TYPES_TO_PERMISSIONS,
} from '@/constants';

import { useCurrentUserPermissions } from '@/hooks/auth';
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
    required: {
      type: Boolean,
      default: false,
    },
    disabled: {
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
    const { checkReadAccess } = useCurrentUserPermissions();

    const items = computed(() =>
      Object.values(TEMPLATE_TESTING_TEST_TYPES)
        .filter(type => checkReadAccess(TEMPLATE_TESTING_TESTS_TYPES_TO_PERMISSIONS[type]))
        .map(type => ({
          value: type,
          text: t(`templateTesting.testTypes.${type}`),
        })),
    );

    return {
      t,
      items,
    };
  },
};
</script>

<style lang="scss" scoped>
.template-testing-test-type-field {
  min-width: 200px;
}
</style>
