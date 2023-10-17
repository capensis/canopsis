<template>
  <c-speed-dial
    v-if="hasCreateAnyViewAccess || hasUpdateAnyViewAccess || hasDeleteAnyViewAccess"
    v-bind="wrapperProps"
    transition="slide-y-reverse-transition"
  >
    <template #activator="{ bind: speedDialBind }">
      <v-tooltip
        :right="tooltipRight"
        :left="tooltipLeft"
        z-index="10"
      >
        <template #activator="{ on: tooltipOn }">
          <v-btn
            class="primary"
            v-on="tooltipOn"
            v-bind="{ ...speedDialBind, ...buttonProps }"
          >
            <v-icon>settings</v-icon>
            <v-icon>close</v-icon>
          </v-btn>
        </template>
        <span>{{ $t('layout.sideBar.buttons.settings') }}</span>
      </v-tooltip>
    </template>
    <v-tooltip
      v-if="hasUpdateAnyViewAccess || hasDeleteAnyViewAccess"
      :right="tooltipRight"
      :left="tooltipLeft"
      z-index="10"
    >
      <template #activator="{ on }">
        <v-btn
          v-on="on"
          :input-value="isNavigationEditingMode"
          color="blue darken-4"
          small
          dark
          fab
          @click.stop="$emit('toggleEditingMode')"
        >
          <v-icon
            dark
            small
          >
            edit
          </v-icon>
          <v-icon
            dark
            small
          >
            done
          </v-icon>
        </v-btn>
      </template>
      <span>{{ $t('layout.sideBar.buttons.edit') }}</span>
    </v-tooltip>
    <v-tooltip
      v-if="hasCreateAnyViewAccess"
      :right="tooltipRight"
      :left="tooltipLeft"
      z-index="10"
    >
      <template #activator="{ on }">
        <v-btn
          v-on="on"
          color="green darken-4"
          small
          dark
          fab
          @click.stop="showCreateViewModal"
        >
          <v-icon
            dark
            small
          >
            add
          </v-icon>
        </v-btn>
      </template>
      <span>{{ $t('layout.sideBar.buttons.create') }}</span>
    </v-tooltip>
  </c-speed-dial>
</template>

<script>
import { MODALS } from '@/constants';

import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';
import layoutNavigationEditingModeMixin from '@/mixins/layout/navigation/editing-mode';

export default {
  mixins: [
    permissionsTechnicalViewMixin,
    layoutNavigationEditingModeMixin,
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
    wrapperProps: {
      type: Object,
      default: () => ({
        direction: 'top',
        bottom: true,
        right: true,
        fixed: true,
      }),
    },
    buttonProps: {
      type: Object,
      default: () => ({
        fab: true,
        dark: true,
      }),
    },
  },
  data() {
    return {
      isVSpeedDialOpen: false,
    };
  },
  methods: {
    showCreateViewModal() {
      this.$modals.show({
        name: MODALS.createView,
      });
    },
  },
};
</script>
