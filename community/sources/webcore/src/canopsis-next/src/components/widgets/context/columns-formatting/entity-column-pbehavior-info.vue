<template>
  <v-tooltip
    v-if="pbehaviorInfo && icon"
    top
  >
    <template #activator="{ on }">
      <v-icon
        color="primary"
        v-on="on"
      >
        {{ icon }}
      </v-icon>
    </template>
    <span>{{ pbehaviorInfo.name }}</span>
  </v-tooltip>
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
