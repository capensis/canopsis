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
      field-filter-editor(v-model="settings.widget.parameters.mfilter", :hiddenFields="['title']")
      v-divider
      field-date-interval(v-model="settings.widget.parameters.dateInterval")
      v-divider
      field-stat-selector(v-model="settings.widget.parameters.stat")
      v-divider
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
import FieldStatSelector from './fields/stats/stat-selector.vue';

export default {
  name: SIDE_BARS.statsParetoSettings,
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldFilterEditor,
    FieldDateInterval,
    FieldStatSelector,
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
