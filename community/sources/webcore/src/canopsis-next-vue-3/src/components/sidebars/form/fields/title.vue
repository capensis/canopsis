<template>
  <widget-settings-item :title="label || $t('common.title')" :optional="!required">
    <c-text-field
      :model-value="modelValue"
      :placeholder="placeholder || $t('settings.widgetTitle')"
      :name="name"
      :required="required"
      :autofocus="autofocus"
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

/**
 * Component to set the title of a widget
 *
 * @prop {String} [modelValue] - Value of the title
 */
export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: String,
      default: '',
    },
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'title',
    },
    label: {
      type: String,
      required: false,
    },
    placeholder: {
      type: String,
      required: false,
    },
    autofocus: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],
};
</script>
