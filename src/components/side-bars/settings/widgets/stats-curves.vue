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
      field-periods-number(v-model="settings.widget.parameters.periods")
      v-divider
      field-stats-select(v-model="settings.widget.parameters.stats")
      v-divider
      field-stats-colors(:stats="settings.widget.parameters.stats", v-model="settings.widget.parameters.statsColors")
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetSettingsMixin from '@/mixins/widget/settings';
import { SIDE_BARS } from '@/constants';

import FieldRowGridSize from '../partial/fields/row-grid-size.vue';
import FieldTitle from '../partial/fields/title.vue';
import FieldDuration from '../partial/fields/duration.vue';
import FieldDateTimeSelect from '../partial/fields/date-time-select.vue';
import FieldPeriodsNumber from '../partial/fields/periods-number.vue';
import FieldStatsSelect from '../partial/fields/stats-select.vue';
import FieldStatsColors from '../partial/fields/stats-colors.vue';

export default {
  name: SIDE_BARS.statsCurvesSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldDuration,
    FieldDateTimeSelect,
    FieldPeriodsNumber,
    FieldStatsSelect,
    FieldStatsColors,
  },
  mixins: [entitiesStatsMixin, widgetSettingsMixin],
  data() {
    const { widget, rowId } = this.config;

    return {
      settings: {
        rowId,
        widget: cloneDeep(widget),
      },
    };
  },
};
</script>

