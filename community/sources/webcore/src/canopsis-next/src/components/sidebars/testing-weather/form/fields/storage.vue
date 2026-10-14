<template>
  <widget-settings-item
    :title="title"
    :disabled="disabled"
    hide-on-disabled
  >
    <c-storage-field
      :model-value="modelValue"
      :required="required"
      :disabled="disabled"
      :name="name"
      @add="$emit('add')"
      @edit="$emit('edit', $event)"
      @remove="$emit('remove')"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: String,
      default: '',
    },
    title: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: true,
    },
    name: {
      type: String,
      default: 'parameters.directory',
    },
  },

  emits: ['update:modelValue', 'add', 'edit', 'remove'],
};
</script>
