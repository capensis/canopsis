<template>
  <v-layout column>
    <field-title v-field="form.title" />
    <field-periodic-refresh v-field="form.parameters" />
    <field-columns
      v-field="form.parameters.widget_columns"
      :template="form.parameters.widget_columns_template"
      :templates="entityColumnsWidgetTemplates"
      :templates-pending="entityColumnsWidgetTemplatesPending"
      :label="$t('settings.columnNames')"
      :type="$constants.ENTITIES_TYPES.entity"
      :excluded-columns="excludedWidgetColumns"
      @update:template="updateWidgetColumnsTemplate"
    />
    <widget-settings-group :title="$t('settings.advancedSettings')">
      <field-columns
        v-field="form.parameters.active_alarms_columns"
        :template="form.parameters.active_alarms_columns_template"
        :templates="alarmColumnsWidgetTemplates"
        :templates-pending="alarmColumnsWidgetTemplatesPending"
        :label="$t('settings.activeAlarmsColumns')"
        :type="$constants.ENTITIES_TYPES.alarm"
        @update:template="updateActiveAlarmsColumnsTemplate"
      />
      <field-columns
        v-field="form.parameters.resolved_alarms_columns"
        :template="form.parameters.resolved_alarms_columns_template"
        :templates="alarmColumnsWidgetTemplates"
        :templates-pending="alarmColumnsWidgetTemplatesPending"
        :label="$t('settings.resolvedAlarmsColumns')"
        :type="$constants.ENTITIES_TYPES.alarm"
        @update:template="updateResolvedAlarmsColumnsTemplate"
      />
      <field-quick-date-interval-type
        v-if="showInterval"
        v-field="form.parameters.default_time_range"
        :ranges="availabilityRanges"
      />
      <field-availability-display-parameter v-field="form.parameters.default_display_parameter" />
      <field-availability-display-show-type v-field="form.parameters.default_show_type" />
      <field-filters
        v-if="showFilter"
        v-field="form.parameters.mainFilter"
        :filters="form.filters"
        :widget-id="widgetId"
        :addable="filterAddable"
        :editable="filterEditable"
        with-entity
        @update:filters="updateField('filters', $event)"
      />
      <export-csv-form
        v-if="showExport"
        v-field="form.parameters.export_settings"
        :type="$constants.ENTITIES_TYPES.entity"
        :templates="entityColumnsWidgetTemplates"
        :templates-pending="entityColumnsWidgetTemplatesPending"
        :excluded-columns="excludedWidgetColumns"
        without-infos-attributes
      />
    </widget-settings-group>
  </v-layout>
</template>

<script>
import { computed } from 'vue';
import { omit } from 'lodash';

import { AVAILABILITY_QUICK_RANGES, CONTEXT_WIDGET_COLUMNS } from '@/constants';

import { formMixin } from '@/mixins/form';

import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldColumns from '@/components/sidebars/form/fields/columns.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import FieldQuickDateIntervalType from '@/components/sidebars/form/fields/quick-date-interval-type.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';
import ExportCsvForm from '@/components/sidebars/form/export-csv.vue';

import FieldAvailabilityDisplayShowType from './fields/availability-display-show-type.vue';
import FieldAvailabilityDisplayParameter from './fields/availability-display-parameter.vue';

export default {
  components: {
    ExportCsvForm,
    FieldAvailabilityDisplayShowType,
    FieldAvailabilityDisplayParameter,
    FieldFilters,
    FieldQuickDateIntervalType,
    WidgetSettingsGroup,
    FieldColumns,
    FieldTitle,
    FieldPeriodicRefresh,
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
    widgetId: {
      type: String,
      required: false,
    },
    entityColumnsWidgetTemplates: {
      type: Array,
      required: false,
    },
    entityColumnsWidgetTemplatesPending: {
      type: Boolean,
      default: false,
    },
    alarmColumnsWidgetTemplates: {
      type: Array,
      required: false,
    },
    alarmColumnsWidgetTemplatesPending: {
      type: Boolean,
      default: false,
    },
    showInterval: {
      type: Boolean,
      default: false,
    },
    showFilter: {
      type: Boolean,
      default: false,
    },
    filterAddable: {
      type: Boolean,
      default: false,
    },
    filterEditable: {
      type: Boolean,
      default: false,
    },
    showExport: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const availabilityRanges = computed(() => AVAILABILITY_QUICK_RANGES);
    const excludedWidgetColumns = computed(() => Object.values(
      omit(CONTEXT_WIDGET_COLUMNS, [
        'id',
        'name',
        'categoryName',
        'type',
        'impactLevel',
        'infos',
        'links',
      ]),
    ));

    const updateWidgetColumnsTemplate = (template, columns) => {
      emit('update:widget-columns-template', template, columns);
    };

    const updateActiveAlarmsColumnsTemplate = (template, columns) => {
      emit('update:active-alarms-columns-template', template, columns);
    };

    const updateResolvedAlarmsColumnsTemplate = (template, columns) => {
      emit('update:resolved-alarms-columns-template', template, columns);
    };

    return {
      availabilityRanges,
      excludedWidgetColumns,
      updateWidgetColumnsTemplate,
      updateActiveAlarmsColumnsTemplate,
      updateResolvedAlarmsColumnsTemplate,
    };
  },
};
</script>
