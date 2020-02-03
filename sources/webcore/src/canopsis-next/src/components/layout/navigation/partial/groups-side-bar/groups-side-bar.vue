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
      v-model="groupss",
      :component-data="{ props: { expand: true, dark: true, focusable: true } }",
      :options="draggableOptions",
      element="v-expansion-panel"
    )
      groups-side-bar-group(
        v-for="group in groupss",
        :key="group._id",
        :group.sync="group",
        :isEditingMode="isEditingMode",
        :draggableOptions="draggableOptions",
      )
    v-divider
    template(v-if="isGroupsOrderChanged")
      v-btn.primary(@click="resetGroups") {{ $t('common.submit') }}
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
      groupss: [],
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
      return this.availableGroups.some((group, index) => this.groupss[index]._id !== group._id ||
        group.views.some((view, viewIndex) => this.groupss[index].views[viewIndex]._id !== view._id));
    },
  },
  watch: {
    availableGroups: {
      deep: true,
      immediate: true,
      handler(groups) {
        this.setGroups(groups);
      },
    },
  },
  methods: {
    setGroups(groups = []) {
      this.groupss = groups.map(group => ({
        ...group,

        views: [...group.views],
      }));
    },

    resetGroups() {
      this.setGroups(this.availableGroups);
    },

    submit() {
      // const original
      const { groups, views } = this.groupss.reduce((acc, group, index) => {
        const originalGroup = this.groups.find(({ _id: id }) => id === group._id);
        const isGroupsOrderChanged = this.availableGroups[index]._id === group._id;

        if (isGroupsOrderChanged) {
          acc.groups.push({ ...group, position: index });
        }

        const views = acc.groups.views.reduce((acc, view, index) => {
          // originalGroup.
        }, []);

        if (views.length) {
          acc.views.concat(views);
        }
      }, { groups: [], views: [] });
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
</style>
