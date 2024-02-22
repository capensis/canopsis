<template>
  <c-speed-dial
    v-if="hasAccessToPrivateView || hasCreateAnyViewAccess || hasUpdateAnyViewAccess || hasDeleteAnyViewAccess"
    v-model="isVSpeedDialOpen"
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
            v-bind="{ ...speedDialBind, ...buttonProps }"
            v-on="tooltipOn"
          >
            <v-icon :small="buttonProps.small">
              settings
            </v-icon>
            <v-icon :small="buttonProps.small">
              close
            </v-icon>
          </v-btn>
        </template>
        <span>{{ $t('layout.sideBar.buttons.settings') }}</span>
      </v-tooltip>
    </template>
    <v-tooltip
      v-if="hasAccessToPrivateView || hasUpdateAnyViewAccess || hasDeleteAnyViewAccess"
      :right="tooltipRight"
      :left="tooltipLeft"
      z-index="10"
    >
      <template #activator="{ on }">
        <v-btn
          :input-value="isNavigationEditingMode"
          color="blue darken-4"
          small
          dark
          fab
          v-on="on"
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
          color="green darken-4"
          small
          dark
          fab
          v-on="on"
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
      <span>{{ $t('layout.sideBar.buttons.createView') }}</span>
    </v-tooltip>
    <v-tooltip
      v-if="hasAccessToPrivateView"
      :right="tooltipRight"
      :left="tooltipLeft"
      z-index="10"
    >
      <template #activator="{ on }">
        <v-btn
          color="blue darken-3"
          small
          dark
          fab
          v-on="on"
          @click.stop="showCreatePrivateViewModal"
        >
          <v-icon
            size="40"
            dark
          >
            $vuetify.icons.person_lock
          </v-icon>
        </v-btn>
      </template>
      <span>{{ $t('layout.sideBar.buttons.createPrivateView') }}</span>
    </v-tooltip>
  </c-speed-dial>
</template>

<script>
import { MODALS } from '@/constants';

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
