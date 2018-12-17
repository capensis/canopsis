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
          v-list-tile
            v-list-tile-title
              router-link(:to="{ name: 'exploitation-event-filter' }")
                v-layout(justify-space-between)
                  span.black--text {{ $t('eventFilter.title') }}
                  v-icon.ml-2 list
      v-menu(bottom, offset-y)
        v-btn.white--text(slot="activator", flat) {{ $t('common.administration') }}
        v-list.pb-0
          v-list-tile
            v-list-tile-title
              router-link(:to="{ name: 'admin-rights' }")
                v-layout(justify-space-between)
                  span.black--text {{ $t('common.rights') }}
                  v-icon verified_user
          v-list-tile
            v-list-tile-title
              router-link(:to="{ name: 'admin-users' }")
                v-layout(justify-space-between)
                  span.black--text {{ $t('common.users') }}
                  v-icon people
          v-list-tile
            v-list-tile-title
              router-link(:to="{ name: 'admin-roles' }")
                v-layout(justify-space-between)
                  span.black--text {{ $t('common.roles') }}
                  v-icon supervised_user_circle
          v-list-tile
            v-list-tile-title
              router-link(:to="{ name: 'admin-parameters' }")
                v-layout(justify-space-between)
                  span.black--text {{ $t('common.parameters') }}
                  v-icon settings
      v-menu(bottom, offset-y)
        v-btn.white--text(slot="activator", flat) {{ currentUser._id }}
        v-list.pb-0
          v-list-tile
            v-list-tile-title
              v-layout
                div {{ $t('user.firstName') }} :
                div.px-1(v-if="currentUser.firstname") {{ currentUser.firstname }}
                div.px-1.font-italic(v-else) {{ $t('common.undefined') }}
          v-divider
          v-list-tile
            v-list-tile-title
              v-layout
                div {{ $t('user.lastName') }} :
                div.px-1(v-if="currentUser.lastname") {{ currentUser.lastname }}
                div.px-1.font-italic(v-else) {{ $t('common.undefined') }}
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
                v-btn.primary(@click.stop="editDefaultView", small, fab, depressed)
                  v-icon edit
          v-divider
          v-list-tile.error.white--text(@click.prevent="logout")
            v-list-tile-title
              v-layout(align-center)
                div {{ $t('common.logout') }}
                v-icon.pl-1.white--text exit_to_app
    template(v-if="isShownGroupsTopBar", slot="extension")
      groups-top-bar
</template>

<script>
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
  computed: {
    defaultViewTitle() {
      const userDefaultView = this.getViewById(this.currentUser.defaultview);
      return userDefaultView ? userDefaultView.title : null;
    },
  },
  methods: {
    editDefaultView() {
      this.showModal({
        name: this.$constants.MODALS.selectView,
        config: {
          action: (viewId) => {
            const user = { ...this.currentUser, defaultview: viewId };
            this.editUserAccount(user);
          },
        },
      });
    },
    async editUserAccount(user) {
      await this.editUser({ user });
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
