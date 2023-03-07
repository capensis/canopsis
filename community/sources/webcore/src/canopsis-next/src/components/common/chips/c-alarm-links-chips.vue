<template lang="pug">
  c-alarm-actions-chips(
    v-if="hasAccessToLinks",
    :items="links",
    :small="small",
    :inline-count="inlineCount",
    item-text="text",
    item-value="url",
    @select="openLink",
    @activate="activate"
  )
    template(#item="{ item }")
      v-icon.mr-1(color="white", small) {{ item.icon }}
      span {{ item.text }}
</template>

<script>
import {
  ALARM_LIST_ACTIONS_TYPES,
  DEFAULT_LINKS_INLINE_COUNT,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
} from '@/constants';

import { harmonizeLinks } from '@/helpers/links';

import { authMixin } from '@/mixins/auth';

export default {
  mixins: [authMixin],
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    inlineCount: {
      type: [Number, String],
      default: DEFAULT_LINKS_INLINE_COUNT,
    },
    small: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    hasAccessToLinks() {
      return this.checkAccess(BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[ALARM_LIST_ACTIONS_TYPES.links]);
    },

    links() {
      return harmonizeLinks(this.alarm.links).map(link => ({
        text: link.label,
        icon: link.icon_name,
        url: link.url,
        color: 'grey',
      }));
    },
  },
  methods: {
    openLink(url) {
      window.open(url, '_blank');
    },

    activate() {
      this.$emit('activate');
    },
  },
};
</script>
