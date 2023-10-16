<template>
  <widget-settings-item
    :title="label || $t('common.title')"
    :optional="!required"
  >
    <v-text-field
      v-field="value"
      v-validate="rules"
      :placeholder="placeholder || $t('settings.widgetTitle')"
      :name="name"
      :error-messages="errors.collect(name)"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

/**
 * Component to set the title of a widget
 *
 * @prop {String} [value] - Value of the title
 */
export default {
  inject: ['$validator'],
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: String,
      required: false,
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
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },
  },
};
</script>
