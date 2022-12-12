<template lang="pug">
  v-tooltip(top)
    template(#activator="{ on }")
      v-icon.c-extra-details__badge.brown.darken-1.white--text(v-on="on", small) {{ icon }}
    div.text-md-center
      strong {{ $t('alarm.actions.iconsTitles.grouping') }}
      v-layout(column)
        div {{ $tc('alarm.actions.iconsFields.rule', rules.length) }}&nbsp;:
        div.rule-name(v-for="rule in rules", :key="rule.id") &nbsp;{{ rule.name }}
      div {{ $t('alarm.actions.iconsFields.parents') }} : {{ total }}
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
