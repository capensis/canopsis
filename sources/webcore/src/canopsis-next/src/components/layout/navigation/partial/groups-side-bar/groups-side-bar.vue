<template lang="pug">
  div
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
        :class="{ ordering: isGroupsOrderChanged }",
        :component-data="{ props: { expand: true, dark: true, focusable: true } }",
        :options="draggableOptions",
        element="v-expansion-panel"
      )
        groups-side-bar-group(
          v-for="(group, index) in mutatedGroups",
          :key="group._id",
          :group.sync="mutatedGroups[index]",
          :isEditingMode="isEditingMode",
          :isGroupsOrderChanged="isGroupsOrderChanged",
          :draggableOptions="draggableOptions"
        )
      v-divider
      v-fade-transition
        div.v-overlay.v-overlay--active(v-show="isGroupsOrderChanged")
          v-btn.primary(@click="submit") {{ $t('common.submit') }}
          v-btn(@click="resetMutatedGroups") {{ $t('common.cancel') }}
      groups-settings-button(
        tooltipRight,
        :isEditingMode="isEditingMode",
        @toggleEditingMode="toggleEditingMode"
      )
    v-fade-transition
      div.v-overlay.v-overlay--active.content-overlay(v-show="isGroupsOrderChanged")
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { groupSchema } from '@/store/schemas';

import entitiesViewMixin from '@/mixins/entities/view';
import layoutNavigationGroupsBarMixin from '@/mixins/layout/navigation/groups-bar';
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
    layoutNavigationGroupsBarMixin,

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
      return this.checkIsGroupsOrderChanged(this.availableGroups, this.mutatedGroups);
    },
  },
  watch: {
    availableGroups: {
      deep: true,
      immediate: true,
      handler(groups, oldGroups) {
        const isGroupsOrderChanged = this.checkIsGroupsOrderChanged(groups, oldGroups);

        if (isGroupsOrderChanged) {
          this.setMutatedGroups(groups);
        }
      },
    },
  },
  methods: {
    checkIsGroupsOrderChanged(groups = [], anotherGroups = []) {
      return this.checkIsEntityOrderChanged(groups, anotherGroups, (entity = {}, anotherEntity = {}) =>
        this.checkIsEntityOrderChanged(entity.views, anotherEntity.views));
    },

    checkIsEntityOrderChanged(entities = [], anotherEntities = [], callback = () => false) {
      return entities.length !== anotherEntities.length ||
        entities.some((entity, index) => {
          const anotherEntity = anotherEntities[index] || {};

          return entity._id !== anotherEntity._id || callback(entity, anotherEntity);
        });
    },

    /**
     * Set mutated groups method
     *
     * @param {Array} [groups=[]] - New mutated groups
     */
    setMutatedGroups(groups = []) {
      this.mutatedGroups = groups.map(group => ({
        ...group,

        views: [...group.views],
      }));
    },

    /**
     * Reset mutated groups method
     */
    resetMutatedGroups() {
      this.setMutatedGroups(this.availableGroups);
    },

    /**
     * Get requests array for group views
     *
     * @param {Array} [views=[]] - Views with updated ordering
     * @param {Array} [originalViews=[]] - Original views from store
     * @returns {Array<Promise>}
     */
    getGroupViewsRequests(views = [], originalViews = []) {
      return views.reduce((viewAcc, view, viewIndex) => {
        const originalView = originalViews[viewIndex];

        if (originalView && originalView._id !== view._id) {
          const viewForUpdate = { ...view, position: viewIndex };

          viewAcc.push(this.updateViewWithoutStore({ data: viewForUpdate, id: view._id }));
        }

        return viewAcc;
      }, []);
    },

    /**
     * Get requests array for groups
     *
     * @param {Array} [groups=[]] - Groups with updated ordering
     * @param {Array} [originalGroups=[]] - Original groups from store
     * @returns {Array<Promise>}
     */
    getGroupsRequests(groups = [], originalGroups = []) {
      return groups.reduce((acc, group, index) => {
        const isGroupsOrderChanged = originalGroups[index]._id !== group._id;
        const originalGroup = originalGroups.find(({ _id: id }) => id === group._id);

        if (isGroupsOrderChanged) {
          const groupForUpdate = { name: group.name, position: index };

          acc.push(this.updateGroup({ data: groupForUpdate, id: group._id }));
        }

        const viewsRequests = this.getGroupViewsRequests(group.views, originalGroup.views);

        if (viewsRequests.length) {
          acc.push(...viewsRequests);
        }

        return acc;
      }, []);
    },

    /**
     * Submit the sidebar ordering
     *
     * @returns {Promise<void>}
     */
    async submit() {
      try {
        const requests = this.getGroupsRequests(this.mutatedGroups, this.availableGroups);

        await Promise.all(requests);
        await this.fetchGroupsList();

        this.$popups.success({ text: this.$t('layout.sideBar.ordering.popups.success') });
      } catch (err) {
        this.$popups.error({ text: this.$t('layout.sideBar.ordering.popups.error') });
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .content-overlay {
    z-index: 3;
  }

  .panel {
    box-shadow: none;

    &.ordering {
      position: absolute;
      z-index: 9;
    }
  }

  .side-bar {
    position: fixed;
    height: 100vh;
    overflow-y: auto;
    z-index: 4;
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
</style>
