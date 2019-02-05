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
      field-stat-end-date-select(
      :title="$t('settings.tstop')",
      v-model="settings.widget.parameters.tstop",
      :duration="settings.widget.parameters.duration"
      )
      v-divider
      field-duration(v-model="settings.widget.parameters.duration", :title="$t('common.step')")
      v-divider
      field-periods-number(v-model="settings.widget.parameters.periods")
      v-divider
      field-stats-select(v-model="settings.widget.parameters.stats")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-stats-colors(
          :stats="settings.widget.parameters.stats",
          v-model="settings.widget.parameters.statsColors"
          )
      v-divider
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetSettingsMixin from '@/mixins/widget/settings';
import { SIDE_BARS } from '@/constants';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldDuration from './fields/stats/duration.vue';
import FieldStatEndDateSelect from './fields/stats/stat-end-date-select.vue';
import FieldPeriodsNumber from './fields/stats/periods-number.vue';
import FieldStatsSelect from './fields/stats/stats-select.vue';
import FieldStatsColors from './fields/stats/stats-colors.vue';

export default {
  name: SIDE_BARS.statsCurvesSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldDuration,
    FieldStatEndDateSelect,
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

