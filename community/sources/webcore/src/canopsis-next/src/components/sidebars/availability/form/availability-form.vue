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
    </widget-settings-group>
  </v-layout>
</template>

<script>
import { AVAILABILITY_QUICK_RANGES } from '@/constants';

import { formMixin } from '@/mixins/form';

import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldColumns from '@/components/sidebars/form/fields/columns.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import FieldQuickDateIntervalType from '@/components/sidebars/form/fields/quick-date-interval-type.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';

import FieldAvailabilityDisplayShowType from './fields/availability-display-show-type.vue';
import FieldAvailabilityDisplayParameter from './fields/availability-display-parameter.vue';

export default {
  components: {
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
  },
  computed: {
    availabilityRanges() {
      return AVAILABILITY_QUICK_RANGES;
    },
  },
  methods: {
    updateWidgetColumnsTemplate(template, columns) {
      this.$emit('update:widget-columns-template', template, columns);
    },

    updateActiveAlarmsColumnsTemplate(template, columns) {
      this.$emit('update:active-alarms-columns-template', template, columns);
    },
  },
};
</script>
