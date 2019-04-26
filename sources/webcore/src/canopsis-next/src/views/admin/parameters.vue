<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.parameters') }}
    v-list
      v-list-tile
        v-list-tile-title {{ $t('parameters.groupsNavigationType.title') }}
        v-list-tile-content
          v-select(
          :items="groupsNavigationOptions",
          :value="groupsNavigationType",
          @change="setGroupsNavigationType"
          )
</template>

<script>
import { GROUPS_NAVIGATION_TYPES } from '@/constants';

import appMixin from '@/mixins/app';

export default {
  mixins: [appMixin],
  data() {
    return {
      groupsNavigationOptions: [
        {
          text: this.$t('parameters.groupsNavigationType.items.sideBar'),
          value: GROUPS_NAVIGATION_TYPES.sideBar,
        },
        {
          text: this.$t('parameters.groupsNavigationType.items.topBar'),
          value: GROUPS_NAVIGATION_TYPES.topBar,
        },
      ],
    };
  },
};
</script>
