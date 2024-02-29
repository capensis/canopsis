<template>
  <widget-settings
    :submitting="submitting"
    divider
    @submit="submit"
  >
    <availability-form
      v-model="form"
      :widget-id="widget._id"
      :entity-columns-widget-templates="entityColumnsWidgetTemplates"
      :entity-columns-widget-templates-pending="widgetTemplatesPending"
      :alarm-columns-widget-templates="alarmColumnsWidgetTemplates"
      :alarm-columns-widget-templates-pending="widgetTemplatesPending"
      :show-interval="hasAccessToInterval"
      :show-filter="hasAccessToListFilters"
      :filter-addable="hasAccessToAddFilter"
      :filter-editable="hasAccessToEditFilter"
      @update:widgetColumnsTemplate="updateWidgetColumnsTemplate"
      @update:activeAlarmsColumnsTemplate="updateActiveAlarmsColumnsTemplate"
    />
  </widget-settings>
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { widgetTemplatesMixin } from '@/mixins/widget/templates';
import { permissionsWidgetsAvailabilityFilters } from '@/mixins/permissions/widgets/availability/filters';
import { permissionsWidgetsAlarmStatisticsInterval } from '@/mixins/permissions/widgets/availability/interval';

import WidgetSettings from '../partials/widget-settings.vue';

import AvailabilityForm from './form/availability-form.vue';

export default {
  name: SIDE_BARS.availabilitySettings,
  components: {
    WidgetSettings,
    AvailabilityForm,
  },
  mixins: [
    widgetSettingsMixin,
    widgetTemplatesMixin,
    permissionsWidgetsAvailabilityFilters,
    permissionsWidgetsAlarmStatisticsInterval,
  ],
  methods: {
    updateWidgetColumnsTemplate(template, columns) {
      this.$set(this.form.parameters, 'widgetColumnsTemplate', template);
      this.$set(this.form.parameters, 'widgetColumns', columns);
    },

    updateActiveAlarmsColumnsTemplate(template, columns) {
      this.$set(this.form.parameters, 'activeAlarmsColumnsTemplate', template);
      this.$set(this.form.parameters, 'activeAlarmsColumns', columns);
    },
  },
};
</script>
