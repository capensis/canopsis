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
      field-stats-selector(v-model="settings.widget.parameters.stats", required)
      v-divider
      field-stats-colors(
        :stats="settings.widget.parameters.stats",
        v-model="settings.widget.parameters.statsColors"
      )
      v-divider
      field-stats-points-styles(
        :stats="settings.widget.parameters.stats",
        v-model="settings.widget.parameters.statsPointsStyles"
      )
      v-divider
      field-stats-annotation-line(v-model="settings.widget.parameters.annotationLine")
      v-divider
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldFilterEditor from '@/components/sidebars/settings/fields/common/filter-editor.vue';
import FieldDateInterval from '@/components/sidebars/settings/fields/stats/date-interval.vue';
import FieldStatsSelector from '@/components/sidebars/settings/fields/stats/stats-selector.vue';
import FieldStatsColors from '@/components/sidebars/settings/fields/stats/stats-colors.vue';
import FieldStatsPointsStyles from '@/components/sidebars/settings/fields/stats/stats-points-style.vue';
import FieldStatsAnnotationLine from '@/components/sidebars/settings/fields/stats/annotation-line.vue';

export default {
  name: SIDE_BARS.statsCurvesSettings,
  components: {
    FieldTitle,
    FieldFilterEditor,
    FieldDateInterval,
    FieldStatsSelector,
    FieldStatsColors,
    FieldStatsPointsStyles,
    FieldStatsAnnotationLine,
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
};
</script>
