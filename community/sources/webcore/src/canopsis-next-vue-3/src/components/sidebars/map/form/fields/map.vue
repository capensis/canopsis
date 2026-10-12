<template>
  <widget-settings-item :title="$t('common.map')">
    <c-map-field
      :value="modelValue"
      name="parameters.map"
      required
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
      default: '',
    },
  },

  emits: ['update:modelValue'],
};
</script>
