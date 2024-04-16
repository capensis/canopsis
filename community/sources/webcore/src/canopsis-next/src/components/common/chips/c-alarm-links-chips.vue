<template>
  <c-links-chips
    v-if="hasAccessToLinks"
    :links="alarm.links"
    :small="small"
    :inline-count="inlineCount"
    :category="category"
    :only-icon="onlyIcon"
    v-on="$listeners"
  />
</template>

<script>
import { ALARM_LIST_ACTIONS_TYPES, BUSINESS_USER_PERMISSIONS_ACTIONS_MAP } from '@/constants';

import { authMixin } from '@/mixins/auth';

export default {
  inject: ['$system'],
  mixins: [authMixin],
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    inlineCount: {
      type: [Number, String],
      required: false,
    },
    small: {
      type: Boolean,
      default: false,
    },
    category: {
      type: String,
      required: false,
    },
    onlyIcon: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    hasAccessToLinks() {
      return this.checkAccess(BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[ALARM_LIST_ACTIONS_TYPES.links]);
    },
  },
};
</script>
