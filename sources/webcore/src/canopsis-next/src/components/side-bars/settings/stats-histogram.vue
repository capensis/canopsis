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
      field-stats-selector(v-model="settings.widget.parameters.stats", required)
      v-divider
      field-stats-colors(
        :stats="settings.widget.parameters.stats",
        v-model="settings.widget.parameters.statsColors"
      )
      v-divider
      field-stats-annotation-line(v-model="settings.widget.parameters.annotationLine")
      v-divider
    v-btn.primary(data-test="submitStatsHistogramButton", @click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import widgetSettingsMixin from '@/mixins/widget/settings';
import { SIDE_BARS } from '@/constants';

import FieldTitle from './fields/common/title.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';
import FieldDateInterval from './fields/stats/date-interval.vue';
import FieldStatsSelector from './fields/stats/stats-selector.vue';
import FieldStatsColors from './fields/stats/stats-colors.vue';
import FieldStatsAnnotationLine from './fields/stats/annotation-line.vue';

export default {
  name: SIDE_BARS.statsHistogramSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldTitle,
    FieldFilterEditor,
    FieldDateInterval,
    FieldStatsSelector,
    FieldStatsColors,
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

