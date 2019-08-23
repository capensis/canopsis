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
      field-date-interval(v-model="settings.widget.parameters.dateInterval", :hiddenFields="['periodValue']")
      v-divider
      field-filter-editor(v-model="settings.widget.parameters.mfilter", :hiddenFields="['title']")
      v-divider
      field-stats-selector(v-model="settings.widget.parameters.stats", required, withTrend, withSorting)
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-default-sort-column(
          v-model="settings.widget.parameters.sort",
          :columns="defaultSortColumns",
          :columnsLabel="$t('settings.columnName')"
          )
          v-divider
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldDateInterval from './fields/stats/date-interval.vue';
import FieldStatsSelector from './fields/stats/stats-selector.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';

export default {
  name: SIDE_BARS.statsTableSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldDateInterval,
    FieldStatsSelector,
    FieldFilterEditor,
    FieldDefaultSortColumn,
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
  computed: {
    defaultSortColumns() {
      return Object.keys(this.settings.widget.parameters.stats).map(key => ({ label: key, value: key }));
    },
  },
};

</script>
