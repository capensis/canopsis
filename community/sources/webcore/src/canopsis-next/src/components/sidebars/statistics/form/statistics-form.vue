<template lang="pug">
  v-layout(column)
    field-title(v-field="form.title")
    field-periodic-refresh(v-field="form.parameters.periodic_refresh")
    field-main-parameter(
      v-field="form.parameters.mainParameter",
      :type="type"
    )
    field-statistics-columns(v-model="form.parameters.widgetColumns")
    widget-settings-group(:title="$t('settings.advancedSettings')")
      field-title(
        v-field="form.parameters.table_title",
        :label="$tc('common.header')",
        :placeholder="$t('settings.headerTitle')"
      )
      field-quick-date-interval-type(v-field="form.parameters.default_time_range")
      field-filters(
        v-field="form.parameters.mainFilter",
        :filters="form.filters",
        :widget-id="widget._id",
        addable,
        editable,
        with-entity,
        @update:filters="updateField('filters', $event)"
      )
</template>

<script>
import { KPI_RATING_SETTINGS_TYPES } from '@/constants';

import { formMixin } from '@/mixins/form';

import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldQuickDateIntervalType from '@/components/sidebars/form/fields/quick-date-interval-type.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';

import FieldMainParameter from './fields/main-parameter.vue';
import FieldStatisticsColumns from './fields/statistics-columns.vue';

export default {
  components: {
    WidgetSettingsGroup,
    FieldTitle,
    FieldPeriodicRefresh,
    FieldQuickDateIntervalType,
    FieldFilters,
    FieldMainParameter,
    FieldStatisticsColumns,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    type: {
      type: Number,
      default: KPI_RATING_SETTINGS_TYPES.entity,
    },
  },
};
</script>
