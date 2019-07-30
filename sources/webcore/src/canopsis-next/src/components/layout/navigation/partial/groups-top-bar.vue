<template lang="pug">
  v-layout.secondary.groups-wrapper
    v-tabs(ref="tabs", color="secondary", show-arrows, dark)
      template(v-if="hasReadAnyViewAccess")
        v-menu(
        v-for="group in availableGroups",
        :key="group._id",
        content-class="group-v-menu-content secondary",
        close-delay="0",
        open-on-hover,
        offset-y,
        bottom,
        dark
        )
          div.v-btn.v-btn--flat.theme--dark(
          :data-test="`dropDownButton-group-${group._id}`",
          slot="activator"
          )
            span {{ group.name }}
            v-btn(
            data-test="editGroupButton",
            v-show="isEditingMode",
            depressed,
            small,
            icon,
            @click.stop="showEditGroupModal(group)"
            )
              v-icon(small) edit
            v-icon(dark) arrow_drop_down
          v-list(:data-test="`dropDownZone-group-${group._id}`")
            v-list-tile(
            v-for="view in group.views",
            :key="view._id",
            :to="getViewLink(view)"
            )
              v-list-tile-title
                span {{ view.title }}
                v-btn.edit-view-button(
                :data-test="`editViewButton-view-${view._id}`",
                v-show="checkViewEditButtonAccessById(view._id)",
                color="grey darken-2",
                depressed,
                small,
                icon,
                @click.prevent="showEditViewModal(view)"
                )
                  v-icon(small) edit
                v-btn.duplicate-view-button(
                :data-test="`copyViewButton-view-${view._id}`",
                v-show="isEditingMode",
                depressed,
                small,
                icon,
                color="grey darken-2",
                @click.prevent="showDuplicateViewModal(view)"
                )
                  v-icon(small) file_copy
    groups-settings-button(
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

import GroupsSettingsButton from './groups-settings-button.vue';

export default {
  components: { GroupsSettingsButton },
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

    .v-menu__activator .v-btn {
      text-transform: none;
    }

    .v-speed-dial--bottom.v-speed-dial--absolute {
      bottom: -10px;
    }

    .v-speed-dial--right.v-speed-dial--absolute {
      right: 25px;
    }
  }

  .group-v-menu-content {
    .v-list {
      background-color: inherit;

      .v-list__tile__title {
        height: 28px;
        line-height: 28px;
      }

      .edit-view-button, .duplicate-view-button {
        vertical-align: top;
        margin: 0 0 0 8px;
      }
    }
  }
</style>
