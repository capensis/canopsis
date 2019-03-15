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
      field-date-interval(v-model="settings.widget.parameters.dateInterval")
      v-divider
      field-filter-editor(v-model="settings.widget.parameters.mfilter")
      v-divider
      field-stats-select(v-model="settings.widget.parameters.stats")
      v-divider
      field-filter-editor(v-model="settings.widget.parameters.mfilter")
      v-divider
      field-stats-colors(
      :stats="settings.widget.parameters.stats",
      v-model="settings.widget.parameters.statsColors"
      )
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';
import FieldDateInterval from './fields/stats/date-interval.vue';
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
    FieldFilterEditor,
    FieldDateInterval,
    FieldStatsSelect,
    FieldStatsColors,
  },
  mixins: [widgetSettingsMixin],
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

