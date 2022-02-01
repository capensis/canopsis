<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      field-date-interval(v-model="settings.widget.parameters.dateInterval")
      v-divider
      field-filter-editor(
        v-model="settings.widget.parameters.mfilter",
        :hidden-fields="['title']",
        :entities-type="$constants.ENTITIES_TYPES.entity"
      )
      v-divider
      field-stat-selector(v-model="settings.widget.parameters.stat")
      v-divider
      field-stats-colors(
        :stats="stats",
        v-model="settings.widget.parameters.statsColors"
      )
      v-divider
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep, omit } from 'lodash';

import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldFilterEditor from '@/components/sidebars/settings/fields/common/filter-editor.vue';
import FieldDateInterval from '@/components/sidebars/settings/fields/stats/date-interval.vue';
import FieldStatSelector from '@/components/sidebars/settings/fields/stats/stat-selector.vue';
import FieldStatsColors from '@/components/sidebars/settings/fields/stats/stats-colors.vue';

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
      const { stat } = this.config.widget.parameters;
      const stats = {
        Accumulation: {},
      };

      stats[stat.title] = omit(stat, ['title']);

      return stats;
    },
  },
};
</script>
