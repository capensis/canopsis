<template lang="pug">
  v-menu(
    :max-width="500",
    :nudge-width="500",
    :close-on-content-click="false",
    origin="left",
    left,
    offset-x,
    offset-overflow
  )
    template(#activator="{ on }")
      v-tooltip(v-on="on", top)
        template(#activator="{ on: tooltipOn }")
          v-btn(v-on="tooltipOn", icon)
            v-icon info
        div {{ $t('calendar.pbehaviorPlanningLegend.title') }}
    v-card
      v-card-text
        template(v-if="exceptionTypes.length")
          div.my-1(v-for="type in exceptionTypes", :key="type._id")
            div.body-1(:style="getStyleForType(type)")
              v-icon.px-1(color="white", small) {{ type.icon_name }}
              strong.ds-ev-title {{ type.name }}
        span.subheading(v-else) {{ $t('calendar.pbehaviorPlanningLegend.noData') }}
</template>

<script>
import { getMostReadableTextColor } from '@/helpers/color';

export default {
  props: {
    exceptionTypes: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    getStyleForType(type = {}) {
      return {
        backgroundColor: type.color,
        color: getMostReadableTextColor(type.color, { level: 'AA', size: 'large' }),
      };
    },
  },
};
</script>
