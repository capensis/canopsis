<template>
  <c-radio-group
    :model-value="modelValue"
    :disabled="disabled"
    color="primary"
    inline
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <template #label>
      {{ t('common.type') }}
    </template>

    <v-radio
      v-for="type in types"
      :key="type"
      :value="type"
      :label="type"
      class="text-uppercase"
    />
  </c-radio-group>
</template>

<script>
import { ROLE_TYPES } from '@/constants';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      default: ROLE_TYPES.ui,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    return {
      t,
      types: Object.values(ROLE_TYPES),
    };
  },
};
</script>
