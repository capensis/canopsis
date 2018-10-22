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
      field-duration(v-model="settings.widget.parameters.duration")
      v-divider
      field-date-select(:title="$t('settings.tstop')", v-model="settings.widget.parameters.tstop")
      v-divider
      field-filter-editor(v-model="settings.widget.parameters.mfilter")
      v-divider
      field-stat-selector(v-model="settings.widget.parameters.stat")
      v-divider
      field-yes-no-mode(v-model="settings.widget.parameters.yesNoMode")
      v-divider
      field-criticity-levels(v-model="settings.widget.parameters.criticityLevels")
      v-divider
      field-levels-colors-selector(v-model="settings.widget.parameters.statColors")
    v-btn(@click="submit", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';
import widgetSettingsMixin from '@/mixins/widget/settings';
import { SIDE_BARS } from '@/constants';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldDuration from './fields/stats/duration.vue';
import FieldDateSelect from './fields/common/date-time-select.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';
import FieldStatSelector from './fields/stats/stat-selector.vue';
import FieldYesNoMode from './fields/stats/yes-no-mode.vue';
import FieldCriticityLevels from './fields/stats/criticity-levels.vue';
import FieldLevelsColorsSelector from './fields/stats/levels-colors-selector.vue';

export default {
  name: SIDE_BARS.statsNumberSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldDuration,
    FieldDateSelect,
    FieldFilterEditor,
    FieldStatSelector,
    FieldYesNoMode,
    FieldCriticityLevels,
    FieldLevelsColorsSelector,
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
