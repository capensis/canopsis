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
        v-btn(slot="activator", flat) {{ currentUser.crecord_name }}
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
                div.px-1(v-if="currentUser.lastname") {{ currentUser.lasttname }}
                div.px-1.font-italic(v-else) {{ $t('common.undefined') }}
          v-divider
          v-list-tile
            v-list-tile-title
              v-layout
                div {{ $t('user.role') }} :
                div.px-1 {{ currentUser.role }}
          v-divider
          v-list-tile
            v-list-tile-title
              v-layout
                div {{ $t('user.defaultView') }} :
                div.px-1(v-if="defaultViewTitle") {{ defaultViewTitle }}
                div.px-1.font-italic(v-else) {{ $t('common.undefined') }}
          v-divider
          v-list-tile.red.darken-4.white--text(@click.prevent="logout")
            v-list-tile-title
              v-layout(align-center)
                div {{ $t('common.logout') }}
                v-icon.pl-1.white--text exit_to_app
</template>

<script>
import find from 'lodash/find';
import forEach from 'lodash/forEach';
import authMixin from '@/mixins/auth';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';

/**
 * Component for the top bar of the application
 *
 * @event toggleSideBar#click
 */
export default {
  mixins: [authMixin, entitiesViewGroupMixin],
  computed: {
    defaultViewTitle() {
      let defaultView = {};
      forEach(this.groups, (group) => {
        defaultView = find(group.views, view => view._id === this.currentUser.defaultview);
        // Return false to exit the loop if view was found
        return !defaultView;
      });

      return defaultView ? defaultView.title : null;
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
