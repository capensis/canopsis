<template lang="pug">
  v-layout.secondary.groups-wrapper
    v-tabs(ref="tabs", color="secondary", show-arrows, dark)
      template(v-if="hasReadAnyViewAccess")
        groups-top-bar-group(
          v-for="group in availableGroups",
          :key="group._id",
          :group="group",
          :isEditingMode="isEditingMode"
        )
    groups-settings-button(
      tooltipLeft,
      :isEditingMode="isEditingMode",
      :wrapperProps="{ direction: 'bottom', absolute: true, right: true, bottom: true }",
      :buttonProps="{ fab: true, dark: true, small: true }",
      @toggleEditingMode="toggleEditingMode"
    )
</template>

<script>
import { groupSchema } from '@/store/schemas';

import vuetifyTabsMixin from '@/mixins/vuetify/tabs';
import entitiesViewGroupMixin from '@/mixins/entities/view/group/index';
import layoutNavigationGroupMenuMixin from '@/mixins/layout/navigation/group-menu';
import registrableMixin from '@/mixins/registrable';

import GroupsSettingsButton from '../groups-settings-button.vue';

import GroupsTopBarGroup from './groups-top-bar-group.vue';

export default {
  components: { GroupsSettingsButton, GroupsTopBarGroup },
  mixins: [
    vuetifyTabsMixin,
    entitiesViewGroupMixin,
    layoutNavigationGroupMenuMixin,

    registrableMixin([groupSchema], 'groups'),
  ],
  watch: {
    groups() {
      this.$nextTick(this.callTabsOnResizeMethod);
    },
  },
  methods: {
    toggleEditingMode() {
      this.isEditingMode = !this.isEditingMode;

      this.$nextTick(this.callTabsOnResizeMethod);
    },
  },
};
</script>

<style lang="scss">
  .groups-wrapper {
    height: 48px;

    .v-speed-dial--bottom.v-speed-dial--absolute {
      bottom: -10px;
    }

    .v-speed-dial--right.v-speed-dial--absolute {
      right: 25px;
    }
  }
</style>
