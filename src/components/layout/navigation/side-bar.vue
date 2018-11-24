<template lang="pug">
  v-navigation-drawer.side-bar.secondary(
  v-model="isOpen",
  :width="$config.SIDE_BAR_WIDTH",
  app,
  )
    div.brand.ma-0.secondary.lighten-1
      v-layout(justify-center, align-center)
        img.my-1(src="@/assets/canopsis.png")
    v-expansion-panel.panel(
    v-if="hasReadAnyViewAccess",
    expand,
    focusable,
    dark
    )
      v-expansion-panel-content(v-for="group in groups", :key="group._id").secondary.white--text
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
              depressed,
              small,
              icon,
              color="grey darken-2",
              @click.prevent="showEditViewModal(view)"
              )
                v-icon(small) edit
    v-divider
    settings-button(
    :isEditingMode="isEditingMode",
    @toggleEditingMode="toggleEditingMode"
    )
</template>

<script>
import layoutNavigationGroupMenuMixin from '@/mixins/layout/navigation/group-menu';

import SettingsButton from './settings-button.vue';

/**
 * Component for the side-bar, on the left of the application
 *
 * @prop {bool} [value=false] - visibility control
 *
 * @event input#update
 */
export default {
  components: { SettingsButton },
  mixins: [layoutNavigationGroupMenuMixin],
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
</style>
