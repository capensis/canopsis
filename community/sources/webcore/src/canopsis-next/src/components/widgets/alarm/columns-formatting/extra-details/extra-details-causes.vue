<template lang="pug">
  div
    v-tooltip.c-extra-details(top)
      v-icon.c-extra-details__badge.brown.darken-1.white--text(
        slot="activator",
        small
      ) {{ icon }}
      div.text-md-center
        strong {{ $t('alarmList.actions.iconsTitles.grouping') }}
        v-layout(row)
          v-flex
            div {{ $tc('alarmList.actions.iconsFields.rule', causesRules.length) }}&nbsp;:
          v-flex
            div(
              v-for="(rule, index) in causesRules",
              :key="rule.id",
              :style="getRuleStyle(index)"
            ) &nbsp;{{ rule.name }}
        div {{ $t('alarmList.actions.iconsFields.causes') }} : {{ causes.total }}
</template>

<script>
import { EVENT_ENTITY_STYLE } from '@/constants';

export default {
  props: {
    causes: {
      type: Object,
      required: true,
    },
  },
  computed: {
    causesRules() {
      return this.causes?.rules ?? [];
    },

    icon() {
      return EVENT_ENTITY_STYLE.groupCauses.icon;
    },
  },
  methods: {
    getRuleStyle(index) {
      if (index % 2 === 1) {
        return { color: '#b5b5b5' };
      }

      return {};
    },
  },
};
</script>
