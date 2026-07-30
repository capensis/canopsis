<template>
  <c-select
    :model-value="value"
    :items="types"
    :label="label || $t('common.type')"
    :name="name"
    :required="required"
    item-title="text"
    item-value="value"
    clearable
    @update:model-value="$emit('update:value', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import { ENTITY_INFO_PROPERTY_TYPES, ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    value: {
      type: Number,
      required: false,
    },
    label: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'type',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:value'],

  setup() {
    const { t } = useI18n();

    const types = computed(() =>
      Object.values(ENTITY_INFO_PROPERTY_TYPES).map(value => ({
        text: t(ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS[value]),
        value,
      })),
    );

    return {
      types,
    };
  },
};
</script>
