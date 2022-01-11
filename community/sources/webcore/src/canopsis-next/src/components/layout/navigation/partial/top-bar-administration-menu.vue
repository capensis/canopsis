<template lang="pug">
  v-menu(v-show="administrationGroupedLinks.length", bottom, offset-y, offset-x)
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
import { sortBy } from 'lodash';

import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';
import entitiesInfoMixin from '@/mixins/entities/info';

import TopBarMenuLink from './top-bar-menu-link.vue';

export default {
  components: { TopBarMenuLink },
  mixins: [authMixin, entitiesInfoMixin],
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
        const links = group.links.filter(this.linksFilterHandler);

        if (links.length) {
          acc.push({ links: sortBy(links, 'title'), title: group.title });
        }

        return acc;
      }, []);
    },

    administrationAccessLinks() {
      return [
        {
          route: { name: 'admin-rights' },
          title: this.$t('common.rights'),
          icon: 'verified_user',
          permission: USERS_PERMISSIONS.technical.permission,
        },
        {
          route: { name: 'admin-roles' },
          title: this.$t('common.roles'),
          icon: 'supervised_user_circle',
          permission: USERS_PERMISSIONS.technical.role,
        },
        {
          route: { name: 'admin-users' },
          title: this.$t('common.users'),
          icon: 'people',
          permission: USERS_PERMISSIONS.technical.user,
        },
      ];
    },

    administrationCommunicationsLinks() {
      return [
        {
          route: { name: 'admin-broadcast-messages' },
          title: this.$t('common.broadcastMessages'),
          icon: '$vuetify.icons.bullhorn',
          permission: USERS_PERMISSIONS.technical.broadcastMessage,
        },
        {
          route: { name: 'admin-playlists' },
          title: this.$t('common.playlists'),
          icon: 'playlist_play',
          permission: USERS_PERMISSIONS.technical.playlist,
        },
      ];
    },

    administrationGeneralLinks() {
      return [
        {
          route: { name: 'admin-engines' },
          title: this.$t('common.engines'),
          icon: '$vuetify.icons.alt_route',
          permission: USERS_PERMISSIONS.technical.engine,
        },
        {
          route: { name: 'admin-parameters' },
          title: this.$t('common.parameters'),
          icon: 'settings',
          permission: USERS_PERMISSIONS.technical.parameters,
        },
        {
          route: { name: 'admin-planning-administration' },
          title: this.$t('common.planning'),
          icon: 'event_note',
          permission: USERS_PERMISSIONS.technical.planning,
        },
        {
          route: { name: 'admin-remediation-administration' },
          title: this.$t('common.remediation'),
          icon: 'assignment',
          permission: USERS_PERMISSIONS.technical.remediation,
        },
      ];
    },

    permissionsWithDefaultType() {
      return [
        USERS_PERMISSIONS.technical.engine,
      ];
    },
  },
  methods: {
    linksFilterHandler({ permission } = {}) {
      return this.permissionsWithDefaultType.includes(permission)
        ? this.checkAccess(permission)
        : this.checkAppInfoAccessByRight(permission) && this.checkReadAccess(permission);
    },
  },
};
</script>

<style lang="scss" scoped>
.top-bar-administration-menu-link /deep/ span {
  margin-left: 8px;
}
</style>
