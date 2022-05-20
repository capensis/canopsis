<template lang="pug">
  div
    v-tooltip.c-extra-details(top)
      v-icon.c-extra-details__badge.brown.darken-1.white--text(
        slot="activator",
        small
      ) {{ icon }}
      div.text-md-center
        strong {{ $t('alarmList.actions.iconsTitles.grouping') }}
        div {{ $t('common.title') }} : {{ ruleName }}
        div {{ $t('alarmList.actions.iconsFields.consequences') }} : {{ total }}
</template>

<script>
import { EVENT_ENTITY_STYLE } from '@/constants';

export default {
  props: {
    total: {
      type: Number,
      required: 0,
    },
    rule: {
      type: Object,
      required: true,
    },
  },
  computed: {
    ruleName() {
      return this.rule?.name ?? '';
    },

    icon() {
      return EVENT_ENTITY_STYLE.groupConsequences.icon;
    },
  },
};
</script>
