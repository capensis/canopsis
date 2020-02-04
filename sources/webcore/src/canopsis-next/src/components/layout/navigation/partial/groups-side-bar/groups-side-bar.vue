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
      v-model="mutatedGroups",
      :class="{ groupAbsolute: isGroupsOrderChanged }",
      :component-data="{ props: { expand: true, dark: true, focusable: true } }",
      :options="draggableOptions",
      element="v-expansion-panel"
    )
      groups-side-bar-group(
        v-for="(group, index) in mutatedGroups",
        :key="group._id",
        :group.sync="mutatedGroups[index]",
        :isEditingMode="issEditingMode",
        :draggableOptions="draggableOptions"
      )
    v-divider
    div.v-overlay.v-overlay--active(v-show="isGroupsOrderChanged")
      v-btn.primary(@click="submit") {{ $t('common.submit') }}
      v-btn(@click="resetGroups") {{ $t('common.cancel') }}
    groups-settings-button(
      tooltipRight,
      :isEditingMode="isEditingMode",
      @toggleEditingMode="toggleEditingMode"
    )
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { groupSchema } from '@/store/schemas';

import entitiesViewMixin from '@/mixins/entities/view';
import layoutNavigationGroupMenuMixin from '@/mixins/layout/navigation/group-menu';
import registrableMixin from '@/mixins/registrable';

import GroupsSettingsButton from '../groups-settings-button.vue';
import AppLogo from '../app-logo.vue';
import AppVersion from '../app-version.vue';
import ActiveSessionsCount from '../active-sessions-count.vue';

import GroupsSideBarGroup from './groups-side-bar-group.vue';

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
    GroupsSideBarGroup,
  },
  mixins: [
    entitiesViewMixin,
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
      mutatedGroups: [],
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

    issEditingMode() {
      return !this.isGroupsOrderChanged && this.isEditingMode;
    },

    draggableOptions() {
      return { animation: VUETIFY_ANIMATION_DELAY, disabled: !this.isEditingMode };
    },

    isGroupsOrderChanged() {
      return this.availableGroups.some((group, index) => this.mutatedGroups[index]._id !== group._id ||
        group.views.some((view, viewIndex) => this.mutatedGroups[index].views[viewIndex]._id !== view._id));
    },
  },
  watch: {
    availableGroups: {
      deep: true,
      immediate: true,
      handler(groups) {
        this.setMutatedGroups(groups);
      },
    },
  },
  methods: {
    setMutatedGroups(groups = []) {
      this.mutatedGroups = groups.map(group => ({
        ...group,

        views: [...group.views],
      }));
    },

    resetGroups() {
      this.setMutatedGroups(this.availableGroups);
    },

    async submit() {
      const promises = this.mutatedGroups.reduce((acc, group, index) => {
        const isGroupsOrderChanged = this.availableGroups[index]._id !== group._id;
        const originalGroup = this.availableGroups.find(({ _id: id }) => id === group._id);

        if (isGroupsOrderChanged) {
          const groupForUpdate = { name: group.name, position: index };

          acc.push(this.updateGroup({ data: groupForUpdate, id: group._id }));
        }

        const viewsPromises = group.views.reduce((viewAcc, view, viewIndex) => {
          const originalView = originalGroup.views[viewIndex];

          if (originalView && originalView._id !== view._id) {
            const viewForUpdate = { ...view, position: viewIndex };

            viewAcc.push(this.updateView({ data: viewForUpdate, id: view._id }));
          }

          return viewAcc;
        }, []);

        if (viewsPromises.length) {
          acc.push(...viewsPromises);
        }

        return acc;
      }, []);

      await Promise.all(promises);

      this.fetchGroupsList();
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

  .logo {
    max-width: 100%;
    max-height: 100%;
    object-fit: scale-down;
  }

  .groupAbsolute {
    position: absolute;
    z-index: 9;
    width: 100%;
  }
</style>
