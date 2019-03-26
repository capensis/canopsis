<template lang="pug">
  v-toolbar.top-bar.primary(
    dense,
    fixed,
    app,
  )
    v-toolbar-side-icon.ml-2.white--text(v-if="isShownGroupsSideBar", @click="$emit('toggleSideBar')")
    v-spacer
    v-toolbar-items
      v-menu(bottom, offset-y)
        v-btn.white--text(slot="activator", flat) {{ $t('common.exploitation') }}
        v-list.pb-0
          v-list-tile(v-for="(link, index) in exploitationLinks", :key="`exploitation-${index}`")
            v-list-tile-title
              router-link(:to="link.route")
                v-layout(justify-space-between)
                  span.black--text {{ link.text }}
                  v-icon.ml-2 {{ link.icon }}
      v-menu(bottom, offset-y)
        v-btn.white--text(slot="activator", flat) {{ $t('common.administration') }}
        v-list.pb-0
          v-list-tile(v-for="(link, index) in administrationLinks", :key="`administration-${index}`")
            v-list-tile-title
              router-link(:to="link.route")
                v-layout(justify-space-between)
                  span.black--text {{ link.text }}
                  v-icon.ml-2 {{ link.icon }}
      v-menu(bottom, offset-y, offset-x)
        v-btn.white--text(slot="activator", flat) {{ currentUser._id }}
        v-list.pb-0
          template(v-if="currentUser.firstname && currentUser.lastname")
            v-list-tile
              v-list-tile-title
                p {{ currentUser.firstname }} {{ currentUser.lastname }}
            v-divider
          v-list-tile
            v-list-tile-title
              v-layout
                div {{ $t('user.role') }} :
                div.px-1 {{ currentUser.role }}
          v-divider
          v-list-tile(two-line)
            v-list-tile-content
              v-layout(align-center)
                v-flex
                  div {{ $t('user.defaultView') }} :
                v-flex(v-if="defaultViewTitle")
                  div.px-1 {{ defaultViewTitle }}
                v-flex(v-else)
                  div.px-1.font-italic {{ $t('common.undefined') }}
                v-btn(@click.stop="editDefaultView", small, fab, icon, depressed)
                  v-icon edit
          v-divider
          v-list-tile
            v-list-tile-title
              v-layout(align-center)
                div {{ $t('common.authKey') }}:
                  span.caption.font-italic  {{ currentUser.authkey }}
          v-divider
          v-list-tile(@click.prevent="logout")
            v-list-tile-title
              v-layout(align-center)
                div.error--text {{ $t('common.logout') }}
                v-icon.pl-1(color="error") exit_to_app
    template(v-if="isShownGroupsTopBar", slot="extension")
      groups-top-bar
</template>

<script>
import { MODALS } from '@/constants';

import appMixin from '@/mixins/app';
import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import entitiesViewMixin from '@/mixins/entities/view/index';
import entitiesUserMixin from '@/mixins/entities/user';

import GroupsTopBar from './groups-top-bar.vue';

/**
 * Component for the top bar of the application
 *
 * @event toggleSideBar#click
 */
export default {
  components: { GroupsTopBar },
  mixins: [appMixin, authMixin, modalMixin, entitiesViewMixin, entitiesUserMixin],
  data() {
    return {
      exploitationLinks: [
        { route: { name: 'exploitation-event-filter' }, text: this.$t('eventFilter.title'), icon: 'list' },
        { route: { name: 'exploitation-pbehaviors' }, text: this.$t('common.pbehaviors'), icon: 'pause' },
        { route: { name: 'exploitation-webhooks' }, text: this.$t('common.webhooks'), icon: '$vuetify.icons.webhook' },
      ],
      administrationLinks: [
        { route: { name: 'admin-rights' }, text: this.$t('common.rights'), icon: 'verified_user' },
        { route: { name: 'admin-users' }, text: this.$t('common.users'), icon: 'people' },
        { route: { name: 'admin-roles' }, text: this.$t('common.roles'), icon: 'supervised_user_circle' },
        { route: { name: 'admin-parameters' }, text: this.$t('common.parameters'), icon: 'settings' },
      ],
    };
  },
  computed: {
    defaultViewTitle() {
      const userDefaultView = this.getViewById(this.currentUser.defaultview);
      return userDefaultView ? userDefaultView.title : null;
    },
  },
  methods: {
    editDefaultView() {
      this.showModal({
        name: MODALS.selectView,
        config: {
          action: (viewId) => {
            const user = { ...this.currentUser, defaultview: viewId };

            return this.editUserAccount(user);
          },
        },
      });
    },
    async editUserAccount(data) {
      await this.createUser({ data });
      await this.fetchCurrentUser();
    },
  },
};
</script>

<style lang="scss" scoped>
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
