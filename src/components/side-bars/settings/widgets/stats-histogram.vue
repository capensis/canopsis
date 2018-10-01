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
      field-date-select(:title="$t('settings.tstop')", v-model="settings.widget.parameters.tstop")
      v-divider
      field-stats-groups(:groups.sync="settings.widget.parameters.groups")
      v-divider
      field-stats-select(v-model="settings.widget.parameters.stats")
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldRowGridSize from '../partial/fields/row-grid-size.vue';
import FieldTitle from '../partial/fields/title.vue';
import FieldStatsGroups from '../partial/fields/stats-groups.vue';
import FieldStatsSelect from '../partial/fields/stats-select.vue';
import FieldDateSelect from '../partial/fields/date-select.vue';

export default {
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldStatsGroups,
    FieldStatsSelect,
    FieldDateSelect,
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

