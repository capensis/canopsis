<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsPointsStyles.title') }}
    v-container(fluid)
      v-layout(v-for="(stat, key) in stats", align-center, :key="key")
        v-flex
          div {{ key }}:
        v-flex
          v-select(:items="pointsStyles", :value="value[key]", @change="updateField(key, $event)")
</template>

<script>
import { STATS_CURVES_POINTS_STYLES } from '@/constants';

import formMixin from '@/mixins/form';

export default {
  mixins: [formMixin],
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

