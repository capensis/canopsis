<template>
  <widget-settings-item :title="$t('settings.defaultSampling')">
    <c-sampling-field
      :value="modelValue"
      :name="name"
      @update:value="$emit('update:modelValue', $event)"
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
      default: undefined,
    },
    name: {
      type: String,
      default: '',
    },
  },
  emits: ['update:modelValue'],
};
</script>
