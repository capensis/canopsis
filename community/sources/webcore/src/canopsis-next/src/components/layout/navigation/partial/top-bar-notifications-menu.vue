<template lang="pug">
  top-bar-menu(:title="$tc('common.notification', 2)", :links="links")
</template>

<script>
import { ROUTE_NAMES, USERS_PERMISSIONS } from '@/constants';

import TopBarMenu from './top-bar-menu.vue';

export default {
  components: { TopBarMenu },
  computed: {
    links() {
      return [
        {
          route: { name: ROUTE_NAMES.notificationInstructionStats },
          icon: 'star_half',
          permission: USERS_PERMISSIONS.technical.notification.instructionStats,
        },
      ];
    },
  },
};
</script>
