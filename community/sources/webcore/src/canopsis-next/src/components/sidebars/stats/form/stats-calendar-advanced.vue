<template>
  <widget-settings-group :title="$t('settings.advancedSettings')">
    <field-filters
      v-field="form.parameters.mainFilter"
      :filters="form.filters"
      :widget-id="widget._id"
      addable
      editable
      with-alarm
      with-entity
      with-pbehavior
      @update:filters="updateField('filters', $event)"
    />
    <field-filters
      :filters="value.filters"
      addable
      editable
      @update:filters="updateField('filters', $event)"
    />
    <v-divider />
    <field-opened-resolved-filter v-field="value.opened" />
    <v-divider />
    <field-switcher
      v-field="value.considerPbehaviors"
      :title="$t('settings.considerPbehaviors.title')"
    />
    <v-divider />
    <field-criticity-levels v-field="value.criticityLevels" />
    <v-divider />
    <field-levels-colors-selector
      v-field="value.criticityLevelsColors"
      color-type="hex"
      hide-suffix
    />
  </widget-settings-group>
</template>

<script>
import { formMixin } from '@/mixins/form';

import FieldOpenedResolvedFilter from '@/components/sidebars/alarm/form/fields/opened-resolved-filter.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';
import FieldSwitcher from '@/components/sidebars/form/fields/switcher.vue';
import FieldCriticityLevels from '@/components/sidebars/stats/form/fields/criticity-levels.vue';
import FieldLevelsColorsSelector from '@/components/sidebars/stats/form/fields/levels-colors-selector.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';

/**
 * Component to regroup the entities list settings fields
 */
export default {
  inject: ['$validator'],
  components: {
    WidgetSettingsGroup,
    FieldOpenedResolvedFilter,
    FieldFilters,
    FieldSwitcher,
    FieldCriticityLevels,
    FieldLevelsColorsSelector,
  },
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    widgetId: {
      type: String,
      required: true,
    },
  },
};
</script>
