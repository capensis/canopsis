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
      field-date-time-select(:title="$t('settings.tstop')", name="tstop", v-model="settings.widget.parameters.tstop")
      v-divider
      field-stats-select(v-model="settings.widget.parameters.stats")
      v-divider
      field-filter-editor(v-model="settings.widget.parameters.mfilter")
    v-btn(@click="submit", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';
import { PAGINATION_LIMIT } from '@/config';
import { SIDE_BARS } from '@/constants';
import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldDuration from './fields/common/duration.vue';
import FieldDateTimeSelect from './fields/common/date-time-select.vue';
import FieldStatsSelect from './fields/stats/stats-select.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';

export default {
  name: SIDE_BARS.statsTableSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldDuration,
    FieldDateTimeSelect,
    FieldStatsSelect,
    FieldFilterEditor,
  },
  mixins: [widgetSettingsMixin],
  data() {
    const { widget, rowId } = this.config;

    return {
      settings: {
        rowId,
        widget: cloneDeep(widget),
        widget_preferences: {
          itemsPerPage: PAGINATION_LIMIT,
        },
      },
    };
  },
};

</script>
