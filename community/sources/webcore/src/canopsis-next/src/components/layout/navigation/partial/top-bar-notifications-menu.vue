<template lang="pug">
  top-bar-menu(:title="$tc('common.notification', 2)", :links="notificationsLinks")
</template>

<script>
import { sortBy } from 'lodash';

import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';
import entitiesInfoMixin from '@/mixins/entities/info';

import TopBarMenu from './top-bar-menu.vue';

export default {
  components: { TopBarMenu },
  mixins: [authMixin, entitiesInfoMixin],
  computed: {
    notificationsLinks() {
      const links = [
        {
          route: { name: 'notification-instruction-stats' },
          title: this.$t('common.instructionRating'),
          icon: 'star_half',
          permission: USERS_PERMISSIONS.technical.notification.instructionStats,
        },
      ];

      const filteredLinks = links.filter(({ permission }) =>
        this.checkAppInfoAccessByPermission(permission) && this.checkReadAccess(permission));

      return sortBy(filteredLinks, 'title');
    },
  },
};
</script>
