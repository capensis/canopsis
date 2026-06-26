<template>
  <v-list class="py-0">
    <template
      v-for="link in links"
      :key="link.key || link.title || link.route?.name"
    >
      <v-list-subheader
        v-if="link.header"
        class="text-subtitle-1"
        @click.stop=""
      >
        {{ link.title }}
      </v-list-subheader>

      <v-divider v-else-if="link.divider" />

      <top-bar-menu-list-link
        v-else
        :link="link"
        @click="handleClick(link)"
        @mouseenter="handleMouseEnter(link, $event)"
      />
    </template>

    <v-menu
      v-if="subItemsShown"
      v-model="subItemsShown"
      :target="subItemsTarget"
      :offset="[0, 4]"
      location="end top"
      origin="start top"
      content-class="top-bar-menu-list-submenu__content"
    >
      <top-bar-menu-list
        :links="parentItem.links"
        :permissions-with-default-type="permissionsWithDefaultType"
        :without-sort="withoutSort"
        @click="handleClick"
      />
    </v-menu>
  </v-list>
</template>

<script>
import { ref, nextTick } from 'vue'

import TopBarMenuListLink from './top-bar-menu-list-link.vue'

export default {
  name: 'TopBarMenuList',

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

  emits: ['click'],

  setup(props, { emit }) {
    const subItemsShown = ref(false)
    const parentItem = ref(undefined)
    const subItemsTarget = ref(undefined)

    const handleMouseEnter = async (link, event) => {
      if (link.links) {
        if (subItemsShown.value) {
          if (link.title === parentItem.value?.title) {
            return
          }

          subItemsShown.value = false

          await nextTick()
        }

        subItemsTarget.value = event.currentTarget
        parentItem.value = link
        subItemsShown.value = true
      } else {
        parentItem.value = undefined
        subItemsTarget.value = undefined
        subItemsShown.value = false
      }
    }

    const handleClick = () => emit('click')

    return {
      subItemsShown,
      parentItem,
      subItemsTarget,
      handleMouseEnter,
      handleClick,
    }
  },
}
</script>
