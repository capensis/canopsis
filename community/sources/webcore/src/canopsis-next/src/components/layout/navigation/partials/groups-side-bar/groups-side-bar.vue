<template>
  <div>
    <v-navigation-drawer
      v-model="isOpen"
      :width="$config.SIDE_BAR_WIDTH"
      :class="{ 'side-bar--editing': isNavigationEditingMode }"
      :ignore-click-outside="isGroupsOrderChanged || hasMaximizedModal"
      class="side-bar secondary"
      app
    >
      <template #prepend>
        <div class="side-bar__brand ma-0 secondary lighten-1">
          <app-logo class="logo" />
          <logged-users-count />
          <app-version class="version" />
        </div>
      </template>
      <section :class="['side-bar__links', { 'side-bar__links--ordering': isGroupsOrderChanged }]">
        <v-layout
          v-if="!mutatedGroups.length && groupsPending"
          class="pa-2"
          justify-center
        >
          <v-progress-circular
            color="primary"
            indeterminate
          />
        </v-layout>
        <c-draggable-list-field
          v-else
          v-model="mutatedGroups"
          :component-data="{ props: expansionPanelsProps }"
          :disabled="!isNavigationEditingMode"
          class="groups-panel"
          draggable=".groups-panel__item--public"
          component="v-expansion-panels"
        >
          <groups-side-bar-group
            v-for="(group, index) in mutatedGroups"
            :key="group._id"
            :group.sync="mutatedGroups[index]"
            :is-groups-order-changed="isGroupsOrderChanged"
            class="groups-panel__item--public"
          />
          <template
            v-if="hasAccessToPrivateView"
            #footer=""
          >
            <groups-side-bar-group
              v-for="privateGroup in privateGroups"
              :key="privateGroup._id"
              :group="privateGroup"
              :is-groups-order-changed="isGroupsOrderChanged"
            />
          </template>
        </c-draggable-list-field>
        <v-divider />
        <groups-side-bar-playlists />
      </section>
      <groups-settings-button
        tooltip-right
        @toggleEditingMode="toggleNavigationEditingMode"
      />
      <v-fade-transition>
        <v-overlay
          :value="isGroupsOrderChanged"
          class="side-bar__overlay"
        >
          <v-btn
            class="primary ma-2"
            @click="submit"
          >
            {{ $t('common.submit') }}
          </v-btn>
          <v-btn
            class="ma-2"
            light
            @click="resetMutatedGroups"
          >
            {{ $t('common.cancel') }}
          </v-btn>
        </v-overlay>
      </v-fade-transition>
    </v-navigation-drawer>
    <v-fade-transition>
      <v-overlay
        :value="isGroupsOrderChanged"
        z-index="8"
      />
    </v-fade-transition>
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { isDeepOrderChanged } from '@/helpers/dragdrop';
import { groupsWithViewsToPositions } from '@/helpers/entities/view/form';

import { entitiesViewMixin } from '@/mixins/entities/view';
import { layoutNavigationGroupsBarMixin } from '@/mixins/layout/navigation/groups-bar';

import GroupsSettingsButton from '../groups-settings-button.vue';
import AppLogo from '../app-logo.vue';
import AppVersion from '../app-version.vue';
import LoggedUsersCount from '../logged-users-count.vue';

import GroupsSideBarGroup from './groups-side-bar-group.vue';
import GroupsSideBarPlaylists from './groups-side-bar-playlists.vue';

const { mapGetters: modalMapGetters } = createNamespacedHelpers('modals');

/**
 * Component for the side-bar, on the left of the application
 *
 * @prop {bool} [value=false] - visibility control
 *
 * @event input#update
 */
export default {
  components: {
    GroupsSettingsButton,
    AppLogo,
    AppVersion,
    LoggedUsersCount,
    GroupsSideBarGroup,
    GroupsSideBarPlaylists,
  },
  mixins: [
    entitiesViewMixin,
    layoutNavigationGroupsBarMixin,
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
    ...modalMapGetters(['hasMaximizedModal']),

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

    expansionPanelsProps() {
      return {
        multiple: true,
        dark: true,
        accordion: true,
        flat: true,
        tile: true,
      };
    },

    privateGroups() {
      return this.availableGroups.filter(group => group.is_private);
    },

    publicGroups() {
      return this.availableGroups.filter(group => !group.is_private);
    },

    isGroupsOrderChanged() {
      return isDeepOrderChanged(
        this.publicGroups,
        this.mutatedGroups,
        '_id',
        (entity = {}, anotherEntity = {}) => isDeepOrderChanged(entity.views, anotherEntity.views),
      );
    },
  },
  watch: {
    publicGroups: {
      deep: true,
      immediate: true,
      handler(groups) {
        this.setMutatedGroups(groups);
      },
    },
  },
  methods: {
    /**
     * Reset mutated groups method
     */
    resetMutatedGroups() {
      this.setMutatedGroups(this.publicGroups);
    },

    /**
     * Set mutated groups method
     *
     * @param {ViewGroupWithViews[]} [groups=[]] - New mutated groups
     */
    setMutatedGroups(groups = []) {
      this.mutatedGroups = groups.map(group => ({
        ...group,

        views: [...group.views],
      }));
    },

    /**
     * Submit the sidebar ordering
     *
     * @returns {Promise<void>}
     */
    async submit() {
      try {
        const data = groupsWithViewsToPositions(this.mutatedGroups);

        await this.updateViewsPositions({ data });
        await this.fetchAllGroupsListWithWidgets();

        this.$popups.success({ text: this.$t('layout.sideBar.ordering.popups.success') });
      } catch (err) {
        this.$popups.error({ text: this.$t('layout.sideBar.ordering.popups.error') });
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .groups-panel {
    position: relative;
    box-shadow: none;
    transition: none;
  }

  .side-bar {
    position: fixed;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    justify-content: stretch;

    &__brand {
      max-height: 48px;
      position: relative;
      display: flex;
      justify-content: center;
      flex-shrink: 0;
      padding: 0.5em 0;

      & ::v-deep .logged-users-count {
        right: 0;
      }
    }

    .v-navigation-drawer__content {
      display: flex;
      flex-direction: column;
      justify-content: stretch;
    }

    &__brand {
      max-height: 48px;
      position: relative;
      display: flex;
      justify-content: center;
      flex-shrink: 0;
      padding: 0.5em 0;

      & ::v-deep .logged-users-count {
        right: 0;
      }
    }

    &__links {
      overflow: auto;
      padding-bottom: 100px;
      width: 100%;

      &--ordering {
        position: absolute;
        z-index: 9;
      }
    }

    &__overlay {
      align-items: flex-start;
      justify-content: flex-start;
    }

    &--editing {
      z-index: 9;

      .groups-panel:after {
        content: '';
        position: absolute;
        top: 100%;
        width: 100%;
        height: 48px;
      }
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
