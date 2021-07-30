<template lang="pug">
  v-menu(v-show="administrationGroupedLinks.length", bottom, offset-y)
    v-btn.white--text(slot="activator", flat) {{ $t('common.administration') }}
    v-list.py-0
      template(v-for="(group, index) in administrationGroupedLinks")
        v-subheader.subheading(:key="`${group.title}-title`", @click.stop="") {{ group.title }}
        top-bar-menu-link.top-bar-administration-menu-link(
          v-for="link in group.links",
          :key="link.title",
          :link="link"
        )
        v-divider(
          v-if="index < administrationGroupedLinks.length - 1",
          :key="`${group.title}-divider`"
        )
</template>

<script>
import { USERS_PERMISSIONS } from '@/constants';

import { layoutNavigationTopBarMenuMixin } from '@/mixins/layout/navigation/top-bar-menu';

import TopBarMenuLink from './top-bar-menu-link.vue';

export default {
  components: { TopBarMenuLink },
  mixins: [layoutNavigationTopBarMenuMixin],
  computed: {
    administrationGroupedLinks() {
      const groupedLinks = [
        {
          title: this.$t('common.access'),
          links: this.administrationAccessLinks,
        },
        {
          title: this.$tc('common.communication', 2),
          links: this.administrationCommunicationsLinks,
        },
        {
          title: this.$t('common.general'),
          links: this.administrationGeneralLinks,
        },
      ];

      return groupedLinks.reduce((acc, group) => {
        const links = this.prepareLinks(group.links);

        if (links.length) {
          acc.push({ links, title: group.title });
        }

        return acc;
      }, []);
    },

    administrationAccessLinks() {
      return [
        {
          route: { name: 'admin-rights' },
          icon: 'verified_user',
          permission: USERS_PERMISSIONS.technical.action,
        },
        {
          route: { name: 'admin-roles' },
          icon: 'supervised_user_circle',
          permission: USERS_PERMISSIONS.technical.role,
        },
        {
          route: { name: 'admin-users' },
          icon: 'people',
          permission: USERS_PERMISSIONS.technical.user,
        },
      ];
    },

    administrationCommunicationsLinks() {
      return [
        {
          route: { name: 'admin-broadcast-messages' },
          icon: '$vuetify.icons.bullhorn',
          permission: USERS_PERMISSIONS.technical.broadcastMessage,
        },
        {
          route: { name: 'admin-playlists' },
          icon: 'playlist_play',
          permission: USERS_PERMISSIONS.technical.playlist,
        },
      ];
    },

    administrationGeneralLinks() {
      return [
        {
          route: { name: 'admin-parameters' },
          icon: 'settings',
          permission: USERS_PERMISSIONS.technical.parameters,
        },
        {
          route: { name: 'admin-planning-administration' },
          icon: 'event_note',
          permission: USERS_PERMISSIONS.technical.planning,
        },
        {
          route: { name: 'admin-remediation-administration' },
          icon: 'assignment',
          permission: USERS_PERMISSIONS.technical.remediation,
        },
      ];
    },
  },
};
</script>

<style lang="scss" scoped>
.top-bar-administration-menu-link /deep/ span {
  margin-left: 8px;
}
</style>
