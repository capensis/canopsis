<template lang="pug">
  div
    v-list.pt-0(expand)
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
    v-btn.primary(
      data-test="paretoDiagramSubmitButton",
      @click="submit"
    ) {{ $t('common.save') }}
</template>

<script>
import { cloneDeep, omit } from 'lodash';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldTitle from './fields/common/title.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';
import FieldDateInterval from './fields/stats/date-interval.vue';
import FieldStatSelector from './fields/stats/stat-selector.vue';
import FieldStatsColors from './fields/stats/stats-colors.vue';

export default {
  name: SIDE_BARS.statsParetoSettings,
  components: {
    FieldTitle,
    FieldFilterEditor,
    FieldDateInterval,
    FieldStatSelector,
    FieldStatsColors,
  },
  mixins: [widgetSettingsMixin],
  data() {
    const { widget } = this.config;

    return {
      settings: {
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
