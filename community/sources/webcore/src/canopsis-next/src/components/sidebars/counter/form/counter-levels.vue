<template>
  <widget-settings-group :title="$t('settings.counterLevels.title')">
    <widget-settings-item :title="$t('settings.counterLevels.fields.counter')">
      <v-select
        v-field="form.counter"
        :label="$t('settings.counterLevels.fields.counter')"
        :items="availableCounters"
        class="select"
        hide-details
        single-line
        dense
      />
    </widget-settings-item>
    <field-criticity-levels v-field="form.values" />
    <field-levels-colors-selector
      v-field="form.colors"
      color-type="hex"
      hide-suffix
    />
  </widget-settings-group>
</template>

<script>
import { AVAILABLE_COUNTERS } from '@/constants';

import FieldCriticityLevels from '@/components/sidebars/stats/form/fields/criticity-levels.vue';
import FieldLevelsColorsSelector from '@/components/sidebars/stats/form/fields/levels-colors-selector.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

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
