<template>
  <c-simple-tooltip v-if="pbehaviorInfo && icon" :content="pbehaviorInfo.name" top>
    <template #activator="{ props: activatorProps }">
      <v-icon v-bind="activatorProps" color="primary">
        {{ icon }}
      </v-icon>
    </template>
  </c-simple-tooltip>
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
