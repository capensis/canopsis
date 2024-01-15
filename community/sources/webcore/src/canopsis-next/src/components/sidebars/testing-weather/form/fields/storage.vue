<template>
  <widget-settings-item
    :title="title"
    :disabled="disabled"
    hide-on-disabled
  >
    <c-storage-field
      v-field="storage"
      v-on="$listeners"
      :required="required"
      :disabled="disabled"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  inject: ['$validator'],
  components: { WidgetSettingsItem },
  model: {
    prop: 'storage',
    event: 'input',
  },
  props: {
    storage: {
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
      default: false,
    },
  },
};
</script>
