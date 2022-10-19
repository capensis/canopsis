<template lang="pug">
  widget-settings-group(:title="$t('settings.counterLevels.title')")
    widget-settings-item(:title="$t('settings.counterLevels.fields.counter')")
      v-select.select(
        v-field="form.counter",
        :label="$t('settings.counterLevels.fields.counter')",
        :items="availableCounters",
        hide-details,
        single-line,
        dense
      )
    v-divider
    field-criticity-levels(v-field="form.values")
    v-divider
    field-levels-colors-selector(
      v-field="form.colors",
      colorType="hex",
      hide-suffix
    )
</template>

<script>
import { AVAILABLE_COUNTERS } from '@/constants';

import FieldCriticityLevels from '@/components/sidebars/settings/fields/stats/criticity-levels.vue';
import FieldLevelsColorsSelector from '@/components/sidebars/settings/fields/stats/levels-colors-selector.vue';
import WidgetSettingsGroup from '@/components/sidebars/settings/partials/widget-settings-group.vue';
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  components: {
    WidgetSettingsItem,
    WidgetSettingsGroup,
    FieldCriticityLevels,
    FieldLevelsColorsSelector,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    availableCounters() {
      return Object.values(AVAILABLE_COUNTERS);
    },
  },
};
</script>
