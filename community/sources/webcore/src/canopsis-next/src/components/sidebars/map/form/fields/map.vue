<template>
  <widget-settings-item :title="$tc('common.map')">
    <c-map-field
      v-field="value"
      name="parameters.map"
    />
  </widget-settings-item>
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  inject: ['$validator'],
  components: { WidgetSettingsItem },
  mixins: [formValidationHeaderMixin],
  props: {
    value: {
      type: String,
      required: false,
    },
  },
};
</script>
