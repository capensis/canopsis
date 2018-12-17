<template lang="pug">
  v-navigation-drawer.side-bar.secondary(
  v-model="isOpen",
  :width="$config.SIDE_BAR_WIDTH",
  disable-resize-watcher,
  app,
  )
    div.brand.ma-0.secondary.lighten-1
      v-layout(justify-center, align-center)
        v-flex.text-xs-center(xs11)
          img.my-1(src="@/assets/canopsis.png")
        v-flex.version.white--text.caption
          div {{ version }}
    v-expansion-panel.panel(
    v-if="hasReadAnyViewAccess",
    expand,
    focusable,
    dark
    )
      v-expansion-panel-content.secondary.white--text(v-for="group in groups", :key="group._id")
        div(slot="header")
          span {{ group.name }}
          v-btn(
          v-show="isEditingMode",
          depressed,
          small,
          icon,
          @click.stop="showEditGroupModal(group)"
          )
            v-icon(small) edit
        v-card.secondary.white--text(v-for="view in getAvailableViewsForGroup(group)", :key="view._id")
          v-card-text
            router-link(:to="{ name: 'view', params: { id: view._id } }")
              span.pl-3 {{ view.title }}
              v-btn(
              v-show="(checkUpdateViewAccessById(view._id) || checkDeleteViewAccessById(view._id)) && isEditingMode",
              color="grey darken-2",
              depressed,
              small,
              icon,
              @click.prevent="showEditViewModal(view)"
              )
                v-icon(small) edit
              v-btn(
              v-show="isEditingMode",
              depressed,
              small,
              icon,
              color="grey darken-2",
              @click.prevent="showDuplicateViewModal(view)"
              )
                v-icon(small) file_copy
    v-divider
    groups-settings-button(
    :isEditingMode="isEditingMode",
    @toggleEditingMode="toggleEditingMode"
    )
</template>

<script>
import versionMixin from '@/mixins/entities/version';
import layoutNavigationGroupMenuMixin from '@/mixins/layout/navigation/group-menu';

import GroupsSettingsButton from './groups-settings-button.vue';

/**
 * Component for the side-bar, on the left of the application
 *
 * @prop {bool} [value=false] - visibility control
 *
 * @event input#update
 */
export default {
  components: { GroupsSettingsButton },
  mixins: [versionMixin, layoutNavigationGroupMenuMixin],
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isOpen: {
      get() {
        return this.value;
      },
      set(value) {
        if (value !== this.value) {
          this.$emit('input', value);
        }
      },
    },
  },
  mounted() {
    this.fetchVersion();
  },
};
</script>

<style scoped>
  a {
    color: inherit;
    text-decoration: none;
  }

  .panel {
    box-shadow: none;
  }

  .side-bar {
    position: fixed;
    height: 100vh;
    overflow-y: auto;
  }

  .brand {
    height: 48px;
    position: relative;
  }

  .version {
    position: absolute;
    bottom: 0;
    right: 0;
    padding-right: 0.5em;
  }
</style>
