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
      field-stat-selector(v-model="settings.widget.parameters.stat")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-stat-display-mode(v-model="settings.widget.parameters.displayMode")
          v-divider
          field-number(v-model="settings.widget.parameters.limit", :title="$t('common.limit')")
          v-divider
          field-sort-order(v-model="settings.widget.parameters.sortOrder")
          v-divider
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import widgetSettingsMixin from '@/mixins/widget/settings';
import { SIDE_BARS } from '@/constants';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldDateInterval from './fields/stats/date-interval.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';
import FieldStatSelector from './fields/stats/stat-selector.vue';
import FieldStatDisplayMode from './fields/stats/stat-display-mode.vue';
import FieldNumber from './fields/common/number.vue';
import FieldSortOrder from './fields/stats/sort-order.vue';

export default {
  name: SIDE_BARS.statsNumberSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldDateInterval,
    FieldFilterEditor,
    FieldStatSelector,
    FieldStatDisplayMode,
    FieldNumber,
    FieldSortOrder,
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
