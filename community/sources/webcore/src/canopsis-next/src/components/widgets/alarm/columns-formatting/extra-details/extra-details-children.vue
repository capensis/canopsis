<template lang="pug">
  div
    v-tooltip.c-extra-details(top, lazy)
      template(#activator="{ on }")
        span.c-extra-details__badge.brown.darken-1(v-on="on")
          v-icon(color="white", small) {{ icon }}
      div.text-md-center
        strong {{ $t('alarm.actions.iconsTitles.grouping') }}
        div {{ $t('common.title') }} : {{ ruleName }}
        div {{ $t('alarm.actions.iconsFields.children') }} : {{ total }}
        div {{ $t('alarm.fields.openedChildren') }} : {{ opened }}
        div {{ $t('alarm.fields.closedChildren') }} : {{ closed }}
</template>

<script>
import { EVENT_ENTITY_ICONS_BY_TYPE } from '@/constants';

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
    opened: {
      type: Number,
      default: 0,
    },
    closed: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    ruleName() {
      return this.rule?.name ?? '';
    },

    icon() {
      return EVENT_ENTITY_ICONS_BY_TYPE.groupChildren;
    },
  },
};
</script>
