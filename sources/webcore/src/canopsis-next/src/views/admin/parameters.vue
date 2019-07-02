<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.parameters') }}
    v-card.ma-2
      v-card-text
        v-list
          v-list-tile
            v-list-tile-title {{ $t('parameters.groupsNavigationType.title') }}
            v-list-tile-content
              v-select.mt-0(
              :items="groupsNavigationOptions",
              :value="groupsNavigationType",
              @change="setGroupsNavigationType"
              )
    v-layout(row, wrap)
      v-flex(xs8, offset-xs2)
        v-card.ma-2
          v-card-text
            user-interface-form
</template>

<script>
import { GROUPS_NAVIGATION_TYPES } from '@/constants';

import appMixin from '@/mixins/app';

import UserInterfaceForm from '@/components/forms/user-interface.vue';

export default {
  components: { UserInterfaceForm },
  mixins: [appMixin],
  computed: {
    groupsNavigationOptions() {
      return [
        {
          text: this.$t('parameters.groupsNavigationType.items.sideBar'),
          value: GROUPS_NAVIGATION_TYPES.sideBar,
        },
        {
          text: this.$t('parameters.groupsNavigationType.items.topBar'),
          value: GROUPS_NAVIGATION_TYPES.topBar,
        },
      ];
    },
  },
};
</script>
