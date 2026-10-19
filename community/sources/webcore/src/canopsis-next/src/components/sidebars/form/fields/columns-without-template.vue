<template>
  <widget-settings-item :title="label" :eager="false">
    <span class="text-body-2 my-2">
      {{ t('common.column', 2) }}
    </span>
    <c-columns-field
      :model-value="modelValue"
      v-bind="$attrs"
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </widget-settings-item>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  inheritAttrs: false,

  props: {
    label: {
      type: String,
      required: true,
    },
    modelValue: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    return {
      t,
    };
  },
};
</script>
