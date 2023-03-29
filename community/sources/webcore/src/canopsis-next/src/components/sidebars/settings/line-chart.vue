<template lang="pug">
  widget-settings(:submitting="submitting", divider, @submit="submit")
    field-title(v-model="form.title")
    field-periodic-refresh(v-model="form.parameters.periodic_refresh")
    widget-settings-group(:title="$t('settings.chart.metricsDisplay')")
      field-alarm-metric-presets(v-model="form.parameters.metrics", with-color)
    widget-settings-group(:title="$t('settings.advancedSettings')")
      field-chart-title(v-model="form.parameters.chart_title")
      field-quick-date-interval-type(v-model="form.parameters.default_time_range")
      field-sampling(v-model="form.parameters.default_sampling")
      field-filters(
        :filters.sync="form.filters",
        addable,
        editable,
        with-entity,
        hide-selector
      )
      field-switcher(v-model="form.parameters.comparison", :title="$t('settings.chart.showComparison')")
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import WidgetSettings from './partials/widget-settings.vue';
import WidgetSettingsGroup from './partials/widget-settings-group.vue';
import FieldTitle from './fields/common/title.vue';
import FieldPeriodicRefresh from './fields/common/periodic-refresh.vue';
import FieldAlarmMetricPresets from './fields/chart/alarm-metric-presets.vue';
import FieldChartTitle from './fields/chart/chart-title.vue';
import FieldQuickDateIntervalType from './fields/common/quick-date-interval-type.vue';
import FieldSampling from './fields/common/sampling.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldSwitcher from './fields/common/switcher.vue';

export default {
  name: SIDE_BARS.lineChartSettings,
  components: {
    WidgetSettings,
    WidgetSettingsGroup,
    FieldTitle,
    FieldPeriodicRefresh,
    FieldAlarmMetricPresets,
    FieldChartTitle,
    FieldQuickDateIntervalType,
    FieldSampling,
    FieldFilters,
    FieldSwitcher,
  },
  mixins: [widgetSettingsMixin],
};
</script>
