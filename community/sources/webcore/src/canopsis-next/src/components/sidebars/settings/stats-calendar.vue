<template lang="pug">
  div
    v-list.pt-0(expand)
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

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { sideBarSettingsWidgetAlarmMixin } from '@/mixins/side-bar/settings/widgets/alarm';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldOpenedResolvedFilter from '@/components/sidebars/settings/fields/alarm/opened-resolved-filter.vue';
import FieldFilters from '@/components/sidebars/settings/fields/common/filters.vue';
import FieldSwitcher from '@/components/sidebars/settings/fields/common/switcher.vue';
import FieldCriticityLevels from '@/components/sidebars/settings/fields/stats/criticity-levels.vue';
import FieldLevelsColorsSelector from '@/components/sidebars/settings/fields/stats/levels-colors-selector.vue';
import AlarmsListModalForm from '@/components/sidebars/settings/forms/alarms-list-modal.vue';
import StatsCalendarAdvancedForm from '@/components/sidebars/settings/forms/stats-calendar-advanced.vue';

/**
 * Component to regroup the entities list settings fields
 */
export default {
  name: SIDE_BARS.statsCalendarSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
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
    const { widget } = this.config;

    return {
      settings: {
        widget: this.prepareWidgetWithAlarmParametersSettings(cloneDeep(widget), true),
      },
    };
  },
  methods: {
    prepareWidgetSettings() {
      const { widget } = this.settings;

      return this.prepareWidgetWithAlarmParametersSettings(widget);
    },
  },
};
</script>
