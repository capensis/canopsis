<template>
  <v-list class="py-0">
    <template v-for="link in links">
      <v-subheader
        v-if="link.header"
        :key="`${link.title}-header`"
        class="text-subtitle-1"
        @click.stop=""
      >
        {{ link.title }}
      </v-subheader>
      <v-divider v-else-if="link.divider" :key="`${link.key}-divider`" />
      <top-bar-menu-list-link
        v-else
        :key="`${link.title}-link`"
        :link="link"
        @mouseenter="handleMouseEnter(link, $event)"
      />
    </template>
    <v-menu
      v-if="subItemsShown"
      v-model="subItemsShown"
      :position-x="subItemsPosition.x"
      :position-y="subItemsPosition.y"
    >
      <top-bar-menu-list
        ref="subItemsList"
        :links="parentItem.links"
        :permissions-with-default-type="permissionsWithDefaultType"
        :without-sort="withoutSort"
      />
    </v-menu>
  </v-list>
</template>

<script>
import { ref, nextTick } from 'vue';

import TopBarMenuListLink from './top-bar-menu-list-link.vue';

const TOP_BAR_MENU_LIST_MIN_WIDTH = 190;

export default {
  name: 'TopBarMenuList', // We need it for recursive
  components: { TopBarMenuListLink },
  props: {
    links: {
      type: Array,
      default: () => [],
    },
    permissionsWithDefaultType: {
      type: Array,
      default: () => [],
    },
    withoutSort: {
      type: Boolean,
      default: true,
    },
  },
  setup() {
    const subItemsList = ref(null);
    const subItemsShown = ref(false);
    const parentItem = ref(undefined);
    const subItemsPosition = ref({ x: 0, y: 0 });

    const handleMouseEnter = async (link, event) => {
      if (link.links) {
        if (subItemsShown.value) {
          if (link.title === parentItem.value?.title) {
            return;
          }

          subItemsShown.value = false;

          await nextTick();
        }

        const { left, top } = event.target.getBoundingClientRect();

        subItemsPosition.value.x = left - TOP_BAR_MENU_LIST_MIN_WIDTH;
        subItemsPosition.value.y = top;
        parentItem.value = link;
        subItemsShown.value = true;

        /**
         * The same logic as in vuetify's VMenu component for activator slot
         */
        await nextTick();

        subItemsList.value?.$el?.offsetWidth;

        setTimeout(() => subItemsPosition.value.x = left - (subItemsList.value?.$el?.offsetWidth ?? 0), 100);
      } else {
        parentItem.value = undefined;
        subItemsShown.value = false;
      }
    };

    return {
      subItemsList,
      subItemsShown,
      parentItem,
      subItemsPosition,
      handleMouseEnter,
    };
  },
};
</script>
