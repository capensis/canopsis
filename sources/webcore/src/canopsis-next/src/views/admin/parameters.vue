<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.parameters') }}
    v-card.ma-2
      v-card-text
        v-list
          v-list-tile
            v-list-tile-title {{ $t('parameters.groupsNavigationType.title') }}
            v-list-tile-content
              v-select(
              :items="groupsNavigationOptions",
              :value="groupsNavigationType",
              @change="setGroupsNavigationType"
              )
    v-card.ma-2
      v-card-text
        v-form
          v-list
            v-list-tile
              v-list-tile-title App title
              v-list-tile-content
                v-text-field(
                )
            v-list-tile
              v-list-tile-title App title
              v-list-tile-content
                v-text-field(
                )
      v-divider
      v-card-actions
        v-btn.primary(type="submit") Submit
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
