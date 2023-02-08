<template lang="pug">
  div
    v-tooltip.c-extra-details(top)
      template(#activator="{ on }")
        span.c-extra-details__badge.brown.darken-1(v-on="on")
          v-icon(color="white", small) {{ icon }}
      div.text-md-center
        strong {{ $t('alarmList.actions.iconsTitles.grouping') }}
        v-layout(column)
          div {{ $tc('alarmList.actions.iconsFields.rule', rules.length) }}&nbsp;:
          div.rule-name(v-for="rule in rules", :key="rule.id") &nbsp;{{ rule.name }}
        div {{ $t('alarmList.actions.iconsFields.parents') }} : {{ total }}
</template>

<script>
import { EVENT_ENTITY_STYLE } from '@/constants';

export default {
  props: {
    rules: {
      type: Array,
      default: () => [],
    },
    total: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    icon() {
      return EVENT_ENTITY_STYLE.groupParents.icon;
    },
  },
};
</script>

<style lang="scss" scoped>
.rule-name:nth-of-type(odd) {
  color: #b5b5b5;
}
</style>
