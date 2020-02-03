<template lang="pug">
  v-navigation-drawer.side-bar.secondary(
    v-model="isOpen",
    :width="$config.SIDE_BAR_WIDTH",
    :class="{ editing: isEditingMode }",
    data-test="groupsSideBar",
    disable-resize-watcher,
    app
  )
    div.brand.ma-0.secondary.lighten-1
      app-logo.logo
      active-sessions-count
      app-version.version
    draggable.panel(
      v-if="hasReadAnyViewAccess",
      element="v-expansion-panel",
      :component-data="{ props: { expand: true, dark: true, focusable: true } }",
      :options="draggableOptions",
      @end="endMoveGroups"
    )
      v-expansion-panel-content.secondary.white--text(
        v-for="group in availableGroups",
        :key="group._id",
        :data-test="`panel-group-${group._id}`"
      )
        div.panel-header(slot="header")
          span(:data-test="`groupsSideBar-group-${group._id}`") {{ group.name }}
          v-btn(
            :data-test="`editGroupButton-group-${group._id}`",
            v-show="isEditingMode",
            depressed,
            small,
            icon,
            @click.stop="showEditGroupModal(group)"
          )
            v-icon(small) edit
        draggable.panel(
          :options="draggableOptions",
          @end=""
        )
          groups-side-bar-item-view-item(
            v-for="view in group.views",
            :key="view._id",
            :view="view",
            :isEditingMode="isEditingMode"
          )
    v-divider
    groups-settings-button(
      tooltipRight,
      :isEditingMode="isEditingMode",
      @toggleEditingMode="toggleEditingMode"
    )
</template>

<script>
import { isEmpty } from 'lodash';
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { groupSchema } from '@/store/schemas';

import layoutNavigationGroupMenuMixin from '@/mixins/layout/navigation/group-menu';
import registrableMixin from '@/mixins/registrable';

import GroupsSettingsButton from './groups-settings-button.vue';
import AppLogo from './app-logo.vue';
import AppVersion from './app-version.vue';
import ActiveSessionsCount from './active-sessions-count.vue';
import GroupsSideBarItemViewItem from './groups-side-bar-item-view-item.vue';

/**
 * Component for the side-bar, on the left of the application
 *
 * @prop {bool} [value=false] - visibility control
 *
 * @event input#update
 */
export default {
  components: {
    Draggable,
    GroupsSettingsButton,
    AppLogo,
    AppVersion,
    ActiveSessionsCount,
    GroupsSideBarItemViewItem,
  },
  mixins: [
    layoutNavigationGroupMenuMixin,

    registrableMixin([groupSchema], 'groups'),
  ],
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      indexesMap: {},
      viewIndexesMap: {},
    };
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

    draggableOptions() {
      return { animation: VUETIFY_ANIMATION_DELAY, disabled: !this.isEditingMode };
    },
    isGroupsOrderChanged() {
      return isEmpty(this.indexesMap);
    },
  },
  methods: {
    endMoveGroups({ newIndex, oldIndex }) {
      const move = this.indexesMap[oldIndex];

      if (move) {
        if (move.originalIndex === newIndex) {
          delete this.indexesMap[oldIndex];
        } else {
          this.indexesMap[newIndex] = {
            ...move,

            oldIndex,
          };
        }
      } else {
        delete this.indexesMap[oldIndex];

        this.indexesMap[newIndex] = {
          oldIndex,
          originalIndex: oldIndex,
        };
      }
    },
  },
};
</script>

<style lang="scss" scoped>
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

    & /deep/ .active-sessions-count {
      position: absolute;
      top: 0;
      right: 0;
    }
  }

  .version {
    position: absolute;
    bottom: 0;
    right: 0;
    padding-right: 0.5em;
    color: white;
    font-size: 0.8em;
    line-height: 1.3em;
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

  .logo {
    max-width: 100%;
    max-height: 100%;
    object-fit: scale-down;
  }
</style>
