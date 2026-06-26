<template>
  <div>
    <v-navigation-drawer
      v-model="isOpen"
      :mobile-breakpoint="mobileBreakpoint"
      :width="sideBarWidth"
      :class="{ 'side-bar--editing': isNavigationEditingMode }"
      class="side-bar"
      color="secondary"
      location="left"
      name="groups-side-bar"
    >
      <template #prepend>
        <div class="side-bar__brand ma-0">
          <app-logo class="logo" />
          <logged-users-count />
          <app-version class="version" />
        </div>
      </template>

      <section :class="['side-bar__links', { 'side-bar__links--ordering': isGroupsOrderChanged }]">
        <c-layout
          v-if="!mutatedGroups.length && groupsPending"
          class="pa-2"
          justify="center"
        >
          <v-progress-circular
            color="primary"
            indeterminate
          />
        </c-layout>

        <c-draggable-list-field
          v-else
          v-model="mutatedGroups"
          :component-data="draggableComponentData"
          :disabled="!isNavigationEditingMode"
          class="groups-panel"
          component="v-expansion-panels"
        >
          <template #item="{ element: group, index }">
            <groups-side-bar-group
              :group="group"
              :is-groups-order-changed="isGroupsOrderChanged"
              class="groups-panel__item--public"
              @update:group="mutatedGroups[index] = $event"
            />
          </template>

          <template
            v-if="hasAccessToPrivateView"
            #footer
          >
            <groups-side-bar-group
              v-for="privateGroup in privateGroups"
              :key="privateGroup._id"
              :group="privateGroup"
              :is-groups-order-changed="isGroupsOrderChanged"
            />
          </template>
        </c-draggable-list-field>

        <v-divider class="side-bar__divider" />

        <groups-side-bar-playlists />
      </section>

      <groups-settings-button
        tooltip-right
        @toggle-editing-mode="toggleNavigationEditingMode"
      />

      <v-overlay
        :model-value="isGroupsOrderChanged"
        class="side-bar__overlay"
        attach=".side-bar"
        persistent
        absolute
      >
        <v-btn
          class="ma-2"
          color="primary"
          variant="elevated"
          @click="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
        <v-btn
          class="ma-2"
          variant="elevated"
          @click="resetMutatedGroups"
        >
          {{ $t('common.cancel') }}
        </v-btn>
      </v-overlay>
    </v-navigation-drawer>

    <v-overlay
      :model-value="isGroupsOrderChanged"
      class="side-bar__overlay"
      persistent
    />
  </div>
</template>

<script>
import { uniq } from 'lodash';
import { mapState } from 'pinia';

import { MEDIA_QUERIES_BREAKPOINTS, SIDE_BAR_WIDTH } from '@/config';

import { isDeepOrderChanged } from '@/helpers/dragdrop';
import { groupsWithViewsToPositions } from '@/helpers/entities/view/form';

import { useModalsStore } from '@/plugins/modals/stores';

import { entitiesViewMixin } from '@/mixins/entities/view';
import { activeViewMixin } from '@/mixins/active-view';
import { layoutNavigationGroupsBarMixin } from '@/mixins/layout/navigation/groups-bar';

import GroupsSettingsButton from '../groups-settings-button.vue';
import AppLogo from '../app-logo.vue';
import AppVersion from '../app-version.vue';
import LoggedUsersCount from '../logged-users-count.vue';

import GroupsSideBarGroup from './groups-side-bar-group.vue';
import GroupsSideBarPlaylists from './groups-side-bar-playlists.vue';

/**
 * Component for the side-bar, on the left of the application
 *
 * @prop {bool} [modelValue=false] - visibility control
 *
 * @event update:modelValue
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
    activeViewMixin,
    layoutNavigationGroupsBarMixin,
  ],

  emits: ['update:modelValue'],

  data() {
    return {
      mutatedGroups: [],
      draggableComponentDataValue: [],
      mobileBreakpoint: MEDIA_QUERIES_BREAKPOINTS.t,
      sideBarWidth: SIDE_BAR_WIDTH,
    };
  },

  computed: {
    ...mapState(useModalsStore, ['hasMaximizedModal']),

    isOpen: {
      get() {
        return this.modelValue;
      },
      set(value) {
        if (value !== this.modelValue) {
          this.$emit('update:modelValue', value);
        }
      },
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

    draggableComponentData() {
      return {
        multiple: true,
        flat: true,
        tile: true,
        modelValue: this.draggableComponentDataValue,
        'onUpdate:modelValue': value => this.draggableComponentDataValue = value,
      };
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

    isOpen(value) {
      if (value) {
        this.calculateDraggableComponentDataValue();
      }
    },
  },

  methods: {
    /**
     * Reset mutated groups method
     */
    resetMutatedGroups() {
      this.setMutatedGroups(this.publicGroups);
    },

    calculateDraggableComponentDataValue() {
      const groupId = this.view?.group?._id;

      if (!groupId) {
        return;
      }

      const index = this.mutatedGroups.findIndex(group => group._id === groupId);

      if (index === -1) {
        return;
      }

      this.draggableComponentDataValue = uniq([...this.draggableComponentDataValue, index]);
    },

    /**
     * Set mutated groups method
     *
     * @param {ViewGroupWithViews[]} [groups=[]] - New mutated groups
     */
    setMutatedGroups(groups = []) {
      this.mutatedGroups = groups.map(group => ({
        ...group,

        views: [...(group.views ?? [])],
      }));

      if (this.isOpen && !this.draggableComponentDataValue.length) {
        this.calculateDraggableComponentDataValue();
      }
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
      } catch {
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

  &__divider {
    border-color: rgba(0, 0, 0, .12) !important;
    opacity: 1;
  }

  &__brand {
    background-color: var(--v-secondary-lighten1);
    max-height: 48px;
    position: relative;
    display: flex;
    justify-content: center;
    flex-shrink: 0;
    padding: 0.5em 0;

    :deep(.logged-users-count) {
      right: 0;
    }
  }

  :deep(.version) {
    position: absolute;
    bottom: 0;
    right: 0;
    padding-right: 0.5em;
    color: white;
    font-size: 0.8em;
    line-height: 1.3em;
  }

  :deep(.v-navigation-drawer__content) {
    display: flex;
    flex-direction: column;
    justify-content: stretch;
  }

  :deep(.v-expansion-panel-title__overlay),
  :deep(.v-card__overlay) {
    opacity: 0 !important;
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
    z-index: 8 !important;
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

.logo {
  max-width: 100%;
  max-height: 100%;
  object-fit: scale-down;
}
</style>
