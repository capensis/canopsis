<template>
  <div v-if="hasSettingsAccess" class="groups-settings-button">
    <c-speed-dial
      :close-on-content-click="false"
      location="top center"
    >
      <template #activator="{ props: activatorProps, isActive }">
        <v-tooltip :location="tooltipLocation">
          <template #activator="{ props: tooltipProps }">
            <v-btn
              v-bind="{ ...activatorProps, ...tooltipProps }"
              :class="{ 'c-speed-dial__btn--active': isActive }"
              class="c-speed-dial__btn"
              color="primary"
              size="x-large"
              icon
            >
              <v-icon icon="settings" />
              <v-icon icon="close" />
            </v-btn>
          </template>
          <span>{{ $t('common.settings') }}</span>
        </v-tooltip>
      </template>

      <div v-if="hasEditAccess" key="edit">
        <v-tooltip :location="tooltipLocation">
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
        </v-tooltip>
      </div>

      <div v-if="hasCreateAnyViewAccess" key="create-view">
        <v-tooltip :location="tooltipLocation">
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
        </v-tooltip>
      </div>

      <div v-if="hasAccessToPrivateView" key="create-private-view">
        <v-tooltip :location="tooltipLocation">
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
        </v-tooltip>
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
</style>
