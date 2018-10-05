<template lang="pug">
  div
    v-list.pt-0(expand)
      field-row-grid-size(
      :rowId.sync="settings.rowId",
      :size.sync="settings.widget.size",
      :availableRows="availableRows",
      @createRow="createRow"
      )
      field-title(v-model="settings.widget.title")
      field-duration(v-model="settings.widget.parameters.duration")
      field-date-select(:title="$t('settings.tstop')", v-model="settings.widget.parameters.tstop")
      field-filter-editor(v-model="settings.widget.parameters.mfilter")
      field-stat-selector(v-model="settings.widget.parameters.stat")
    v-btn(@click="submit", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';
import widgetSettingsMixin from '@/mixins/widget/settings';
import { SIDE_BARS } from '@/constants';

import FieldRowGridSize from '../partial/fields/row-grid-size.vue';
import FieldTitle from '../partial/fields/title.vue';
import FieldDuration from '../partial/fields/duration.vue';
import FieldDateSelect from '../partial/fields/date-time-select.vue';
import FieldFilterEditor from '../partial/fields/filter-editor.vue';
import FieldStatSelector from '../partial/fields/stat-selector.vue';

export default {
  name: SIDE_BARS.statsNumberSettings,
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldDuration,
    FieldDateSelect,
    FieldFilterEditor,
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
