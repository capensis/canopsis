<template lang="pug">
  v-speed-dial(
    v-if="hasAccessToPrivateView || hasCreateAnyViewAccess || hasUpdateAnyViewAccess || hasDeleteAnyViewAccess",
    v-model="isVSpeedDialOpen",
    v-bind="wrapperProps",
    transition="slide-y-reverse-transition"
  )
    template(#activator="")
      v-tooltip(
        :right="tooltipRight",
        :left="tooltipLeft",
        z-index="10",
        custom-activator
      )
        template(#activator="{ on }")
          v-btn.primary(
            v-on="on",
            v-bind="buttonProps",
            :input-value="isVSpeedDialOpen"
          )
            v-icon settings
            v-icon close
        span {{ $t('layout.sideBar.buttons.settings') }}
    v-tooltip(
      v-if="hasAccessToPrivateView || hasUpdateAnyViewAccess || hasDeleteAnyViewAccess",
      :right="tooltipRight",
      :left="tooltipLeft",
      z-index="10",
      custom-activator
    )
      template(#activator="{ on }")
        v-btn(
          v-on="on",
          :input-value="isNavigationEditingMode",
          color="blue darken-4",
          small,
          dark,
          fab,
          @click.stop="$emit('toggleEditingMode')"
        )
          v-icon(dark) edit
          v-icon(dark) done
      span {{ $t('layout.sideBar.buttons.edit') }}
    v-tooltip(
      v-if="hasCreateAnyViewAccess",
      :right="tooltipRight",
      :left="tooltipLeft",
      z-index="10",
      custom-activator
    )
      template(#activator="{ on }")
        v-btn(
          v-on="on",
          color="green darken-4",
          small,
          dark,
          fab,
          @click.stop="showCreateViewModal"
        )
          v-icon(dark) add
      span {{ $t('layout.sideBar.buttons.createView') }}
    v-tooltip(
      v-if="hasAccessToPrivateView",
      :right="tooltipRight",
      :left="tooltipLeft",
      z-index="10",
      custom-activator
    )
      template(#activator="{ on }")
        v-btn(
          v-on="on",
          color="blue darken-3",
          small,
          dark,
          fab,
          @click.stop="showCreatePrivateViewModal"
        )
          v-icon(dark) $vuetify.icons.person_lock
      span {{ $t('layout.sideBar.buttons.createPrivateView') }}
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
