<template>
  <div
    v-if="hasSettingsAccess"
    class="groups-settings-button"
  >
    <v-menu
      v-model="isOpen"
      :close-on-content-click="false"
      location="top center"
      transition="slide-y-reverse-transition"
    >
      <template #activator="{ props: activatorProps }">
        <v-tooltip :location="tooltipLocation">
          <template #activator="{ props: tooltipProps }">
            <v-btn
              color="primary"
              size="x-large"
              icon
              v-bind="{ ...activatorProps, ...tooltipProps }"
            >
              <c-icon v-if="isOpen">close</c-icon>
              <c-icon v-else>settings</c-icon>
            </v-btn>
          </template>
          <span>{{ $t('common.settings') }}</span>
        </v-tooltip>
      </template>

      <c-layout
        class="groups-settings-button__list"
        align="center"
        gap="2"
        column
      >

        <v-tooltip
          v-if="hasAccessToPrivateView"
          :location="tooltipLocation"
        >
          <template #activator="{ props: tooltipProps }">
            <v-btn
              color="blue-darken-3"
              icon
              v-bind="tooltipProps"
              @click.stop="showCreatePrivateViewModal"
            >
              <c-icon size="40">$vuetify.icons.person_lock</c-icon>
            </v-btn>
          </template>
          <span>{{ $t('layout.sideBar.buttons.createPrivateView') }}</span>
        </v-tooltip>

        <v-tooltip
          v-if="hasCreateAnyViewAccess"
          :location="tooltipLocation"
        >
          <template #activator="{ props: tooltipProps }">
            <v-btn
              color="green-darken-4"
              icon
              v-bind="tooltipProps"
              @click.stop="showCreateViewModal"
            >
              <c-icon>add</c-icon>
            </v-btn>
          </template>
          <span>{{ $t('layout.sideBar.buttons.createView') }}</span>
        </v-tooltip>

        <v-tooltip
          v-if="hasEditAccess"
          :location="tooltipLocation"
        >
          <template #activator="{ props: tooltipProps }">
            <v-btn
              :active="isNavigationEditingMode"
              color="blue-darken-4"
              icon
              v-bind="tooltipProps"
              @click.stop="$emit('toggleEditingMode')"
            >
              <c-icon v-if="isNavigationEditingMode">done</c-icon>
              <c-icon v-else>edit</c-icon>
            </v-btn>
          </template>
          <span>{{ $t('layout.sideBar.buttons.edit') }}</span>
        </v-tooltip>
      </c-layout>
    </v-menu>
  </div>
</template>

<script>
import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';
import { layoutNavigationEditingModeMixin } from '@/mixins/layout/navigation/editing-mode';
import { entitiesViewMixin } from '@/mixins/entities/view';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';

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

  data() {
    return {
      isOpen: false,
    };
  },

  computed: {
    hasEditAccess() {
      return this.hasAccessToPrivateView || this.hasUpdateAnyViewAccess || this.hasDeleteAnyViewAccess;
    },

    hasSettingsAccess() {
      return this.hasEditAccess || this.hasCreateAnyViewAccess;
    },

    tooltipLocation() {
      return {
        [this.tooltipLeft]: 'left',
        [this.tooltipRight]: 'right',
      }.true || 'right';
    },
  },

  methods: {
    showCreateViewModal() {
      // TODO: Enable after view modals are migrated.
    },

    showCreatePrivateViewModal() {
      // TODO: Enable after private view modals are migrated.
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

  &__list {
    padding: 8px;
    align-items: center;
  }
}
</style>
