<template lang="pug">
  v-navigation-drawer.side-bar.secondary(
  v-model="isOpen",
  :width="$config.SIDE_BAR_WIDTH",
  :class="{ editing: isEditingMode }"
  data-test="groupsSideBar"
  disable-resize-watcher,
  app,
  )
    div.brand.ma-0.secondary.lighten-1
      img.logo(:src="appLogo")
      div.version {{ version }}
    v-expansion-panel.panel(
    v-if="hasReadAnyViewAccess",
    expand,
    focusable,
    dark
    )
      v-expansion-panel-content.secondary.white--text(
      v-for="group in availableGroups",
      :key="group._id",
      :data-test="`panel-groupName-${group.name}`"
      )
        div.panel-header(slot="header")
          span(:data-test="`groupsSideBar-group-${group._id}`") {{ group.name }}
          v-btn(
          :data-test="`editGroupButton-groupName-${group.name}`",
          v-show="isEditingMode",
          depressed,
          small,
          icon,
          @click.stop="showEditGroupModal(group)"
          )
            v-icon(small) edit
        v-card(
        v-for="view in group.views",
        :key="view._id",
        :color="getColor(view._id)",
        )
          router-link.panel-item-content-link(
          :data-test="`linkView-viewTitle-${view.title}`"
          :title="view.title",
          :to="getViewLink(view)",
          )
            v-card-text.panel-item-content
              v-layout(align-center, justify-space-between)
                v-flex
                  v-layout(align-center)
                    span.pl-2 {{ view.title }}
                v-flex
                  v-layout(justify-end)
                    v-btn.ma-0(
                    :data-test="`editViewButton-viewTitle-${view.title}`",
                    v-show="checkViewEditButtonAccessById(view._id)",
                    depressed,
                    small,
                    icon,
                    @click.prevent="showEditViewModal(view)"
                    )
                      v-icon(small) edit
                    v-btn.ma-0(
                    :data-test="`copyViewButton-viewTitle-${view.title}`",
                    v-show="isEditingMode",
                    depressed,
                    small,
                    icon,
                    @click.prevent="showDuplicateViewModal(view)"
                    )
                      v-icon(small) file_copy
          v-divider
    v-divider
    groups-settings-button(
    :isEditingMode="isEditingMode",
    @toggleEditingMode="toggleEditingMode"
    )
</template>

<script>
import { groupSchema } from '@/store/schemas';

import entitiesInfoMixin from '@/mixins/entities/info';
import layoutNavigationGroupMenuMixin from '@/mixins/layout/navigation/group-menu';
import registrableMixin from '@/mixins/registrable';

import logo from '@/assets/canopsis.png';

import GroupsSettingsButton from './groups-settings-button.vue';

/**
 * Component for the side-bar, on the left of the application
 *
 * @prop {bool} [value=false] - visibility control
 *
 * @event input#update
 */
export default {
  components: { GroupsSettingsButton },
  mixins: [
    entitiesInfoMixin,
    layoutNavigationGroupMenuMixin,

    registrableMixin([groupSchema], 'groups'),
  ],
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isOpen: {
      get() {
        return this.value;
      },
      set(value) {
        if (value !== this.value) {
          this.$emit('input', value);
        }
      },
    },

    isViewActive() {
      return viewId => this.$route.params.id && this.$route.params.id === viewId;
    },

    getColor() {
      return id => (this.isViewActive(id) ? 'secondary white--text lighten-3' : 'secondary white--text lighten-1');
    },

    appLogo() {
      if (this.logo) {
        return this.logo;
      }

      return logo;
    },
  },
};
</script>

<style lang="scss" scoped>
  a {
    color: inherit;
    text-decoration: none;
  }

  .panel {
    box-shadow: none;
  }

  .side-bar {
    position: fixed;
    height: 100vh;
    overflow-y: auto;

    & /deep/ .v-expansion-panel__header {
      height: 48px;
    }
  }

  .brand {
    max-height: 48px;
    position: relative;
    display: flex;
    justify-content: center;
    padding: 0.5em 0;
  }

  .version {
    position: absolute;
    bottom: 0;
    right: 0;
    padding-right: 0.5em;
    color: white;
    font-size: 0.8em;
  }

  .panel-header {
    max-width: 88%;

    span {
      max-width: 100%;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      display: inline-block;
      vertical-align: middle;

      .editing & {
        max-width: 73%;
      }
    }
  }

  .panel-item-content {
    display: -webkit-box;
    display: -ms-flexbox;
    display: flex;
    cursor: pointer;
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    position: relative;
    padding: 12px 24px;
    height: 48px;

    & > div {
      max-width: 100%;
    }

    & /deep/ .v-btn:not(:last-child) {
      margin-right: 0;
    }

    .panel-item-content-link {
      max-width: 100%;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      display: inline-block;
      vertical-align: middle;
    }
  }

  .logo {
    max-width: 100%;
    max-height: 100%;
    object-fit: scale-down;
  }
</style>
