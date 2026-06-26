<template>
  <v-menu
    v-if="preparedLinks.length"
    v-model="openedMenu"
    origin="top center"
    v-bind="$attrs"
  >
    <template #activator="{ props: activatorProps }">
      <v-btn
        class="top-bar-menu__button text-white"
        variant="text"
        v-bind="activatorProps"
      >
        <slot name="title">
          {{ title }}
        </slot>
      </v-btn>
    </template>

    <top-bar-menu-list
      :links="preparedLinks"
      :permissions-with-default-type="permissionsWithDefaultType"
      :without-sort="withoutSort"
      @click="handleClick"
    />
  </v-menu>
</template>

<script>
import { computed, ref, toRef } from 'vue'

import { useTopBarMenu } from './hooks/top-bar-menu'
import TopBarMenuList from './top-bar-menu-list.vue'

export default {
  components: { TopBarMenuList },

  inheritAttrs: false,

  props: {
    title: {
      type: String,
      default: '',
    },
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
      default: false,
    },
  },

  setup(props) {
    const openedMenu = ref(false)
    const { prepareLinks } = useTopBarMenu({
      withoutSort: toRef(props, 'withoutSort'),
      permissionsWithDefaultType: toRef(props, 'permissionsWithDefaultType'),
    })

    const preparedLinks = computed(() => prepareLinks(props.links))
    const handleClick = () => openedMenu.value = false

    return {
      openedMenu,
      preparedLinks,
      handleClick,
    }
  },
}
</script>

<style lang="scss" scoped>
.top-bar-menu__button {
  --v-btn-height: 48px;

  border-radius: 0;
  font-size: 14px;
  font-weight: 500;
  height: 48px;
  letter-spacing: 0;
  line-height: normal;
  min-height: 0;
  padding: 0 16px;
  text-transform: uppercase;

  :deep(.v-btn__content) {
    height: 100%;
  }
}
</style>
