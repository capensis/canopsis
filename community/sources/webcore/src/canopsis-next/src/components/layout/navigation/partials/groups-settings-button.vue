<template>
  <div v-if="hasSettingsAccess" class="groups-settings-button">
    <c-speed-dial
      :close-on-content-click="false"
      v-bind="speedDialProps"
    >
      <template #activator="{ props: activatorProps, isActive }">
        <c-tooltip :location="tooltipLocation">
          <template #activator="{ props: tooltipProps }">
            <v-btn
              v-bind="{ ...activatorProps, ...tooltipProps, ...speedDialBtnProps }"
              :class="{ 'c-speed-dial__btn--active': isActive }"
              class="c-speed-dial__btn"
              icon
            >
              <v-icon icon="settings" />
              <v-icon icon="close" />
            </v-btn>
          </template>
          <span>{{ $t('common.settings') }}</span>
        </c-tooltip>
      </template>

      <div v-if="hasEditAccess" key="edit" class="groups-settings-button__action--edit">
        <c-tooltip :location="tooltipLocation">
          <template #activator="{ props: tooltipProps }">
            <v-btn
              :active="isNavigationEditingMode"
              color="blue-darken-4"
              size="default"
              icon
              v-bind="tooltipProps"
              @click.stop="$emit('toggleEditingMode')"
            >
              <c-icon v-if="isNavigationEditingMode">
                done
              </c-icon>
              <c-icon v-else>
                edit
              </c-icon>
            </v-btn>
          </template>
          <span>{{ $t('layout.sideBar.buttons.edit') }}</span>
        </c-tooltip>
      </div>

      <div
        v-if="hasCreateAnyViewAccess"
        key="create-view"
        class="groups-settings-button__action--create-view"
      >
        <c-tooltip :location="tooltipLocation">
          <template #activator="{ props: tooltipProps }">
            <v-btn
              color="green-darken-4"
              size="default"
              icon
              v-bind="tooltipProps"
              @click.stop="showCreateViewModal"
            >
              <c-icon>add</c-icon>
            </v-btn>
          </template>
          <span>{{ $t('layout.sideBar.buttons.createView') }}</span>
        </c-tooltip>
      </div>

      <div
        v-if="hasAccessToPrivateView"
        key="create-private-view"
        class="groups-settings-button__action--create-private-view"
      >
        <c-tooltip :location="tooltipLocation">
          <template #activator="{ props: tooltipProps }">
            <v-btn
              color="blue-darken-3"
              size="default"
              icon
              v-bind="tooltipProps"
              @click.stop="showCreatePrivateViewModal"
            >
              <c-icon size="40">
                $vuetify.icons.person_lock
              </c-icon>
            </v-btn>
          </template>
          <span>{{ $t('layout.sideBar.buttons.createPrivateView') }}</span>
        </c-tooltip>
      </div>
    </c-speed-dial>
  </div>
</template>

<script>
import { MODALS } from '@/constants';

import { entitiesViewMixin } from '@/mixins/entities/view';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { layoutNavigationEditingModeMixin } from '@/mixins/layout/navigation/editing-mode';
import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';

export default {
  mixins: [
    permissionsTechnicalViewMixin,
    layoutNavigationEditingModeMixin,
    entitiesViewMixin,
    entitiesViewGroupMixin,
  ],

  props: {
    topBar: {
      type: Boolean,
      default: false,
    },
    tooltipRight: {
      type: Boolean,
      default: false,
    },
    tooltipLeft: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['toggleEditingMode'],

  computed: {
    hasEditAccess() {
      return (
        this.hasAccessToPrivateView || this.hasUpdateAnyViewAccess || this.hasDeleteAnyViewAccess
      );
    },

    hasSettingsAccess() {
      return this.hasEditAccess || this.hasCreateAnyViewAccess;
    },

    tooltipLocation() {
      return (
        {
          [this.tooltipLeft]: 'left',
          [this.tooltipRight]: 'right',
        }.true || 'right'
      );
    },

    speedDialProps() {
      if (this.topBar) {
        return {
          contentClass: 'groups-settings-button__content--top-bar',
          location: 'bottom center',
          transition: 'slide-y-transition',
        };
      }

      return {
        location: 'top center',
      };
    },

    speedDialBtnProps() {
      if (this.topBar) {
        return {
          color: 'white',
          size: 'default',
        };
      }

      return {
        color: 'primary',
        size: 'x-large',
      };
    },
  },

  methods: {
    async createViewModalCallback(data) {
      await this.createViewWithPopup({ data });

      return this.fetchAllGroupsListWithWidgetsWithCurrentUser();
    },

    showCreateViewModal() {
      this.$modals.show({
        name: MODALS.createView,
        config: {
          action: this.createViewModalCallback,
        },
      });
    },

    showCreatePrivateViewModal() {
      this.$modals.show({
        name: MODALS.createView,
        config: {
          view: {
            is_private: true,
          },
          title: this.$t('modals.view.create.privateTitle'),
          action: this.createViewModalCallback,
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.groups-settings-button {
  position: fixed;
  right: 16px;
  bottom: 16px;
  z-index: 8;
}

:global(.groups-settings-button__content--top-bar) {
  .groups-settings-button__action--edit {
    order: 3;
  }

  .groups-settings-button__action--create-view {
    order: 2;
  }

  .groups-settings-button__action--create-private-view {
    order: 1;
  }
}

</style>
