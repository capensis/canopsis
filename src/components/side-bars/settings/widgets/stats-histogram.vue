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
      field-title(v-model="settings.widget.title")
      v-divider
      field-duration(v-model="settings.widget.parameters.duration")
      v-divider
      field-date-time-select(:title="$t('settings.tstop')", v-model="settings.widget.parameters.tstop")
      v-divider
      field-stats-groups(v-model="settings.widget.parameters.groups")
      v-divider
      field-stats-select(v-model="settings.widget.parameters.stats")
      v-divider
      field-stats-colors(:stats="settings.widget.parameters.stats", v-model="settings.widget.parameters.statsColors")
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

import widgetSettingsMixin from '@/mixins/widget/settings';
import { SIDE_BARS } from '@/constants';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldDuration from './fields/common/duration.vue';
import FieldStatsGroups from './fields/stats/stats-groups.vue';
import FieldStatsSelect from './fields/stats/stats-select.vue';
import FieldStatsColors from './fields/stats/stats-colors.vue';
import FieldDateTimeSelect from './fields/common/date-time-select.vue';

export default {
  name: SIDE_BARS.statsHistogramSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldDuration,
    FieldStatsGroups,
    FieldStatsSelect,
    FieldStatsColors,
    FieldDateTimeSelect,
  },
  mixins: [widgetSettingsMixin],
  data() {
    const { widget, rowId } = this.config;

    return {
      settings: {
        rowId,
        widget: cloneDeep(widget),
        widget_preferences: {
        },
      },
    };
  },
};
</script>

