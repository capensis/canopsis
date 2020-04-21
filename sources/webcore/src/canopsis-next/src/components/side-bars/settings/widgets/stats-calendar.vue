<template lang="pug">
  div
    v-list.pt-0(expand)
      field-row-grid-size(
        :rowId.sync="settings.rowId",
        :size.sync="settings.widget.size",
        :availableRows="availableRows",
        @createRow="createRow"
      )
      v-divider
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      alarms-list-modal-form(v-model="settings.widget.parameters.alarmsList")
      v-divider
      stats-calendar-advanced-form(v-model="settings.widget.parameters")
      v-divider
    v-btn.primary(data-test="submitStatsCalendarButton", @click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';
import sideBarSettingsWidgetAlarmMixin from '@/mixins/side-bar/settings/widgets/alarm';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldOpenedResolvedFilter from './fields/alarm/opened-resolved-filter.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldSwitcher from './fields/common/switcher.vue';
import FieldCriticityLevels from './fields/stats/criticity-levels.vue';
import FieldLevelsColorsSelector from './fields/stats/levels-colors-selector.vue';
import AlarmsListModalForm from './forms/alarms-list-modal.vue';
import StatsCalendarAdvancedForm from './forms/stats-calendar-advanced.vue';

/**
 * Component to regroup the entities list settings fields
 */
export default {
  name: SIDE_BARS.statsCalendarSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldOpenedResolvedFilter,
    FieldFilters,
    FieldSwitcher,
    FieldCriticityLevels,
    FieldLevelsColorsSelector,
    AlarmsListModalForm,
    StatsCalendarAdvancedForm,
  },
  mixins: [widgetSettingsMixin, sideBarSettingsWidgetAlarmMixin],
  data() {
    const { widget, rowId } = this.config;

    return {
      settings: {
        rowId,
        widget: this.prepareWidgetWithAlarmParametersSettings(cloneDeep(widget), true),
      },
    };
  },
  methods: {
    prepareWidgetSettings() {
      const { widget } = this.settings;

      return this.prepareWidgetWithAlarmParametersSettings(widget);
    },

    prepareWidgetQuery(newQuery, oldQuery) {
      return {
        tstart: oldQuery.tstart,
        tstop: oldQuery.tstop,

        ...newQuery,
      };
    },
  },
};
</script>
