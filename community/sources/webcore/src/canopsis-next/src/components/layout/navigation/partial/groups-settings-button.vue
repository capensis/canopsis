<template lang="pug">
  v-speed-dial(
    data-test="settingsWrapper",
    v-if="hasCreateAnyViewAccess || hasUpdateAnyViewAccess || hasDeleteAnyViewAccess",
    v-model="isVSpeedDialOpen",
    transition="slide-y-reverse-transition",
    v-bind="wrapperProps"
  )
    v-tooltip(
      slot="activator",
      :right="tooltipRight",
      :left="tooltipLeft",
      z-index="10"
    )
      v-btn.primary(
        data-test="settingsViewButton",
        slot="activator",
        :input-value="isVSpeedDialOpen",
        v-bind="buttonProps"
      )
        v-icon settings
        v-icon close
      span {{ $t('layout.sideBar.buttons.settings') }}
    v-tooltip(
      v-if="hasUpdateAnyViewAccess || hasDeleteAnyViewAccess",
      :right="tooltipRight",
      :left="tooltipLeft",
      z-index="10"
    )
      v-btn(
        data-test="editModeButton",
        slot="activator",
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
      z-index="10"
    )
      v-btn(
        data-test="addViewButton",
        slot="activator",
        color="green darken-4",
        small,
        dark,
        fab,
        @click.stop="showCreateViewModal"
      )
        v-icon(dark) add
      span {{ $t('layout.sideBar.buttons.create') }}
</template>

<script>
import { MODALS } from '@/constants';

import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';
import layoutNavigationEditingModeMixin from '@/mixins/layout/navigation/editing-mode';

export default {
  mixins: [
    rightsTechnicalViewMixin,
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
