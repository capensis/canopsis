<template lang="pug">
  v-tooltip(v-if="pbehaviorInfo && icon", top)
    template(#activator="{ on }")
      v-icon(v-on="on", color="primary") {{ icon }}
    span {{ pbehaviorInfo.name }}
</template>

<script>
import { WEATHER_ICONS } from '@/constants';

export default {
  props: {
    pbehaviorInfo: {
      type: Object,
      required: false,
    },
  },
  computed: {
    icon() {
      return WEATHER_ICONS[this.pbehaviorInfo?.canonical_type];
    },
  },
};
</script>
