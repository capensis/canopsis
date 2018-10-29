<template lang="pug">
  v-toolbar.top-bar.white(
    dense,
    fixed,
    clipped-left,
    app,
  )
    div.brand.ma-0.green.darken-4(v-show="$options.filters.mq($mq, { t: true })")
      img(src="@/assets/canopsis.png")
    v-toolbar-side-icon(@click="$emit('toggleSideBar')")
    v-spacer
    v-toolbar-items
      v-menu(bottom, offset-y)
        v-btn(slot="activator", flat) {{ currentUser._id }}
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
                v-btn.green.darken-4.white--text(@click.stop="editDefaultView", small, fab, depressed)
                  v-icon edit
          v-divider
          v-list-tile.red.darken-4.white--text(@click.prevent="logout")
            v-list-tile-title
              v-layout(align-center)
                div {{ $t('common.logout') }}
                v-icon.pl-1.white--text exit_to_app
</template>

<script>
import authMixin from '@/mixins/auth';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesUserMixin from '@/mixins/entities/user';
import modalMixin from '@/mixins/modal/modal';

/**
 * Component for the top bar of the application
 *
 * @event toggleSideBar#click
 */
export default {
  mixins: [authMixin, entitiesViewMixin, entitiesUserMixin, modalMixin],
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

  .top-bar /deep/ .v-toolbar__content {
    padding: 0;
  }
</style>
