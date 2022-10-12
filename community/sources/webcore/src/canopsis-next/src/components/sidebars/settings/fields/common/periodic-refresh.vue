<template lang="pug">
  widget-settings-item(:title="$t('settings.periodicRefresh')", optional)
    periodic-refresh-field(v-field="value", :name="name")
</template>

<script>
import { TIME_UNITS } from '@/constants';

import { formMixin } from '@/mixins/form';

import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  inject: ['$validator'],
  components: { WidgetSettingsItem, PeriodicRefreshField },
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      required: false,
    },
  },
  created() {
    if (!this.value?.unit) {
      this.updateField('unit', TIME_UNITS.second);
    }
  },
};
</script>
