<template lang="pug">
  v-list-group(data-test="statsPointsStyle")
    v-list-tile(slot="activator") {{ $t('settings.statsPointsStyles.title') }}
    v-container(fluid)
      v-layout(v-for="(stat, key) in stats", align-center, :key="key")
        v-flex
          div {{ key }}:
        v-flex(:data-test="`statsPointsStyleSelect-${key}`")
          v-select(
            v-field="value[key]",
            :items="pointsStyles"
          )
</template>

<script>
import { STATS_CURVES_POINTS_STYLES } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    stats: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    pointsStyles() {
      return Object.values(STATS_CURVES_POINTS_STYLES);
    },
  },
};
</script>

