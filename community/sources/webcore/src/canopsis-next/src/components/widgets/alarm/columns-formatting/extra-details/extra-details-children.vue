<template lang="pug">
  v-tooltip(top)
    template(#activator="{ on }")
      v-icon.c-extra-details__badge.brown.darken-1.white--text(v-on="on", small) {{ icon }}
    div.text-md-center
      strong {{ $t('alarm.actions.iconsTitles.grouping') }}
      div {{ $t('common.title') }} : {{ ruleName }}
      div {{ $t('alarm.actions.iconsFields.children') }} : {{ total }}
</template>

<script>
import { EVENT_ENTITY_STYLE } from '@/constants';

export default {
  props: {
    rule: {
      type: Object,
      required: true,
    },
    total: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    ruleName() {
      return this.rule?.name ?? '';
    },

    icon() {
      return EVENT_ENTITY_STYLE.groupChildren.icon;
    },
  },
};
</script>
