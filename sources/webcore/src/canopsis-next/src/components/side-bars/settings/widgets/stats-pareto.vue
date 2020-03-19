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
      field-filter-editor(
        data-test="widgetFilterEditor",
        v-model="settings.widget.parameters.mfilter",
        :hiddenFields="['title']",
        :entitiesType="$constants.ENTITIES_TYPES.entity"
      )
      v-divider
      field-stat-selector(v-model="settings.widget.parameters.stat")
      v-divider
      field-stats-colors(
        :stats="stats",
        v-model="settings.widget.parameters.statsColors"
      )
      v-divider
    copy-widget-id(:widgetId="settings.widget._id")
    v-btn.primary(
      data-test="paretoDiagramSubmitButton",
      @click="submit"
    ) {{ $t('common.save') }}
</template>

<script>
import { cloneDeep, omit } from 'lodash';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';
import CopyWidgetId from '@/components/side-bars/settings/widgets/fields/common/copy-widget-id.vue';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';
import FieldDateInterval from './fields/stats/date-interval.vue';
import FieldStatSelector from './fields/stats/stat-selector.vue';
import FieldStatsColors from './fields/stats/stats-colors.vue';

export default {
  name: SIDE_BARS.statsParetoSettings,
  components: {
    CopyWidgetId,
    FieldRowGridSize,
    FieldTitle,
    FieldFilterEditor,
    FieldDateInterval,
    FieldStatSelector,
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
  computed: {
    stats() {
      const stats = {
        Accumulation: {},
      };

      stats[this.widget.parameters.stat.title] = omit(this.widget.parameters.stat, ['title']);

      return stats;
    },
  },
};
</script>
