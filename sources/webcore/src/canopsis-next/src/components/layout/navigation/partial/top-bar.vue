<template lang="pug">
  v-toolbar.top-bar.primary(
    dense,
    fixed,
    app,
  )
    v-toolbar-side-icon.ml-2.white--text(
    v-if="isShownGroupsSideBar",
    data-test="groupsSideBarButton",
    @click="$emit('toggleSideBar')"
    )
    v-layout.topBarBrand(v-else, fill-height, align-center)
      img.canopsisLogo(src="@/assets/canopsis.png")
      v-layout.version.ml-1(fill-height, align-end)
        div {{ version }}
    v-toolbar-title.white--text.font-weight-regular(v-if="appTitle") {{ appTitle }}
    v-spacer
    v-toolbar-items
      v-menu(v-show="exploitationLinks.length", bottom, offset-y)
        v-btn.white--text(slot="activator", flat) {{ $t('common.exploitation') }}
        v-list.pb-0
          v-list-tile(v-for="(link, index) in exploitationLinks", :key="`exploitation-${index}`")
            v-list-tile-title
              router-link(:to="link.route")
                v-layout(justify-space-between)
                  span.black--text {{ link.text }}
                  v-icon.ml-2 {{ link.icon }}
      v-menu(v-show="administrationLinks.length", bottom, offset-y)
        v-btn.white--text(slot="activator", flat) {{ $t('common.administration') }}
        v-list.pb-0
          v-list-tile(v-for="(link, index) in administrationLinks", :key="`administration-${index}`")
            v-list-tile-title
              router-link(:to="link.route")
                v-layout(justify-space-between)
                  span.black--text {{ link.text }}
                  v-icon.ml-2 {{ link.icon }}
      v-menu(bottom, offset-y, offset-x)
        v-btn.white--text(slot="activator", data-test="userTopBarDropdownButton", flat) {{ currentUser._id }}
        v-list.pb-0
          v-list-tile
            v-list-tile-content
              v-btn.ma-0.pa-1(data-test="userProfileButton", flat, @click.prevent="showEditUserModal")
                v-layout(align-center)
                  v-icon person
                  div.ml-2 {{ $t('user.seeProfile') }}
          v-list-tile
            v-list-tile-content
              v-btn.ma-0.pa-1.error--text(data-test="logoutButton", flat, @click.prevent="logout")
                v-layout(align-center)
                  v-icon exit_to_app
                  div.ml-2 {{ $t('common.logout') }}
    template(v-if="isShownGroupsTopBar", slot="extension")
      groups-top-bar
</template>

<script>
import sha1 from 'sha1';
import { omit, cloneDeep } from 'lodash';

import { MODALS, USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import entitiesUserMixin from '@/mixins/entities/user';
import entitiesInfoMixin from '@/mixins/entities/info';

import GroupsTopBar from './groups-top-bar.vue';

/**
 * Component for the top bar of the application
 *
 * @event toggleSideBar#click
 */
export default {
  components: { GroupsTopBar },
  mixins: [
    authMixin,
    modalMixin,
    entitiesUserMixin,
    entitiesInfoMixin,
  ],
  computed: {
    exploitationLinks() {
      const links = [
        {
          route: { name: 'exploitation-event-filter' },
          text: this.$t('eventFilter.title'),
          icon: 'list',
          right: USERS_RIGHTS.technical.exploitation.eventFilter,
        },
        {
          route: { name: 'exploitation-pbehaviors' },
          text: this.$t('common.pbehaviors'),
          icon: 'pause',
          right: USERS_RIGHTS.technical.exploitation.pbehavior,
        },
        {
          route: { name: 'exploitation-webhooks' },
          text: this.$t('common.webhooks'),
          icon: '$vuetify.icons.webhook',
          right: USERS_RIGHTS.technical.exploitation.webhook,
        },
        {
          route: { name: 'exploitation-snmp-rules' },
          text: this.$t('snmpRules.title'),
          icon: 'assignment',
          right: USERS_RIGHTS.technical.exploitation.snmpRule,
        },
        {
          route: { name: 'exploitation-heartbeats' },
          text: this.$t('heartbeat.title'),
          icon: 'assignment',
          right: USERS_RIGHTS.technical.exploitation.heartbeat,
        },
      ];

      return links.filter(({ right }) => this.checkReadAccess(right));
    },

    administrationLinks() {
      const links = [
        {
          route: { name: 'admin-rights' },
          text: this.$t('common.rights'),
          icon: 'verified_user',
          right: USERS_RIGHTS.technical.action,
        },
        {
          route: { name: 'admin-users' },
          text: this.$t('common.users'),
          icon: 'people',
          right: USERS_RIGHTS.technical.user,
        },
        {
          route: { name: 'admin-roles' },
          text: this.$t('common.roles'),
          icon: 'supervised_user_circle',
          right: USERS_RIGHTS.technical.role,
        },
        {
          route: { name: 'admin-parameters' },
          text: this.$t('common.parameters'),
          icon: 'settings',
          right: USERS_RIGHTS.technical.parameters,
        },
      ];

      return links.filter(({ right }) => this.checkReadAccess(right));
    },
  },
  methods: {
    showEditUserModal() {
      this.showModal({
        name: MODALS.createUser,
        config: {
          title: this.$t('common.profile'),
          user: this.currentUser,
          onlyUserPrefs: true,
          action: async (data) => {
            const editedUser = cloneDeep(this.currentUser);

            if (data.password && data.password !== '') {
              editedUser.shadowpasswd = sha1(data.password);
            }

            await this.createUser({ data: { ...editedUser, ...omit(data, ['password']) } });

            await this.fetchCurrentUser();
          },
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .canopsisLogo {
    max-height: 80%;
    margin-left: 1em;
  }

  .version {
    color: white;
    font-size: 0.7em;
  }

  a {
    text-decoration: none;
    color: inherit;
  }

  .brand {
    display: flex;
    align-items: center;
    margin: 0;
    width: 250px;
    height: 100%;

    img {
      margin: auto;
    }
  }

  .top-bar {
    & /deep/ .v-toolbar__content {
      padding: 0;
    }

    & /deep/ .v-toolbar__extension {
      padding: 0;
    }
  }
</style>
