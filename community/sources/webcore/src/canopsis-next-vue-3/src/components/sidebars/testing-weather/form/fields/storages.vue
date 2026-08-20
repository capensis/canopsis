<template>
  <widget-settings-item
    :title="title"
    optional
  >
    <c-storages-field
      :model-value="modelValue"
      :label="$t('settings.serverStorage')"
      :disabled="disabled"
      :help-text="helpText"
      @add="$emit('add')"
      @edit="(...args) => $emit('edit', ...args)"
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    helpText: {
      type: String,
      required: false,
    },
    title: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue', 'add', 'edit'],
};
</script>
