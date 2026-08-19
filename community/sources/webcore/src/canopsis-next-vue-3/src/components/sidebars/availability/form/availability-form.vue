<template>
  <c-layout column>
    <field-title v-model="localForm.title" />
    <field-periodic-refresh v-model="localForm.parameters.periodic_refresh" />
    <field-columns
      v-model="localForm.parameters.widget_columns"
      :template="localForm.parameters.widget_columns_template"
      :templates="entityColumnsWidgetTemplates"
      :templates-pending="entityColumnsWidgetTemplatesPending"
      :label="$t('settings.columnNames')"
      :type="$constants.ENTITIES_TYPES.entity"
      :excluded-columns="excludedWidgetColumns"
      @update:template="updateWidgetColumnsTemplate"
    />
    <widget-settings-group :title="$t('settings.advancedSettings')">
      <field-columns
        v-model="localForm.parameters.active_alarms_columns"
        :template="localForm.parameters.active_alarms_columns_template"
        :templates="alarmColumnsWidgetTemplates"
        :templates-pending="alarmColumnsWidgetTemplatesPending"
        :label="$t('settings.activeAlarmsColumns')"
        :type="$constants.ENTITIES_TYPES.alarm"
        @update:template="updateActiveAlarmsColumnsTemplate"
      />
      <field-columns
        v-model="localForm.parameters.resolved_alarms_columns"
        :template="localForm.parameters.resolved_alarms_columns_template"
        :templates="alarmColumnsWidgetTemplates"
        :templates-pending="alarmColumnsWidgetTemplatesPending"
        :label="$t('settings.resolvedAlarmsColumns')"
        :type="$constants.ENTITIES_TYPES.alarm"
        @update:template="updateResolvedAlarmsColumnsTemplate"
      />
      <field-quick-date-interval-type
        v-if="showInterval"
        v-model="localForm.parameters.default_time_range"
        :ranges="availabilityRanges"
      />
      <field-availability-display-parameter v-model="localForm.parameters.default_display_parameter" />
      <field-availability-display-show-type v-model="localForm.parameters.default_show_type" />
      <field-filters
        v-if="showFilter"
        v-model="localForm.parameters.mainFilter"
        :filters="localForm.filters"
        :widget-id="widgetId"
        addable
        editable
        with-entity
        @update:filters="localForm.filters = $event"
      />
      <export-csv-form
        v-if="showExport"
        v-model="localForm.parameters.export_settings"
        :type="$constants.ENTITIES_TYPES.entity"
        :templates="entityColumnsWidgetTemplates"
        :templates-pending="entityColumnsWidgetTemplatesPending"
        :excluded-columns="excludedWidgetColumns"
        without-infos-attributes
      />
    </widget-settings-group>
  </c-layout>
</template>

<script>
import { omit } from 'lodash';

import { AVAILABILITY_QUICK_RANGES, CONTEXT_WIDGET_COLUMNS } from '@/constants';

import FieldQuickDateIntervalType from '@/components/sidebars/chart/form/fields/quick-date-interval-type.vue';
import ExportCsvForm from '@/components/sidebars/form/export-csv.vue';
import FieldColumns from '@/components/sidebars/form/fields/columns.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';

import FieldAvailabilityDisplayParameter from './fields/availability-display-parameter.vue';
import FieldAvailabilityDisplayShowType from './fields/availability-display-show-type.vue';

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

  props: {
    modelValue: {
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
    showExport: {
      type: Boolean,
      default: false,
    },
  },

  emits: [
    'update:modelValue',
    'update:widget-columns-template',
    'update:active-alarms-columns-template',
    'update:resolved-alarms-columns-template',
  ],

  computed: {
    localForm: {
      get() {
        return this.modelValue;
      },
      set(value) {
        this.$emit('update:modelValue', value);
      },
    },

    availabilityRanges() {
      return AVAILABILITY_QUICK_RANGES;
    },

    excludedWidgetColumns() {
      return Object.values(
        omit(CONTEXT_WIDGET_COLUMNS, [
          'id',
          'name',
          'component',
          'connector',
          'type',
          'categoryName',
          'impactLevel',
          'enabled',
          'infos',
          'componentInfos',
          'links',
        ]),
      );
    },
  },

  methods: {
    updateWidgetColumnsTemplate(template, columns) {
      this.$emit('update:widget-columns-template', template, columns);
    },

    updateActiveAlarmsColumnsTemplate(template, columns) {
      this.$emit('update:active-alarms-columns-template', template, columns);
    },

    updateResolvedAlarmsColumnsTemplate(template, columns) {
      this.$emit('update:resolved-alarms-columns-template', template, columns);
    },
  },
};
</script>
