<template>
  <v-menu
    v-if="preparedLinks.length"
    v-model="openedMenu"
    v-bind="$attrs"
    offset-y
  >
    <template #activator="{ on }">
      <v-btn
        class="white--text"
        text
        v-on="on"
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
import { computed, ref, toRef } from 'vue';

import { useTopBarMenu } from './hooks/top-bar-menu';
import TopBarMenuList from './top-bar-menu-list.vue';

export default {
  components: { TopBarMenuList },
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Boolean,
      default: false,
    },
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
    const openedMenu = ref(false);

    const { prepareLinks } = useTopBarMenu({
      withoutSort: toRef(props, 'withoutSort'),
      permissionsWithDefaultType: toRef(props, 'permissionsWithDefaultType'),
    });

    const preparedLinks = computed(() => prepareLinks(props.links));

    const handleClick = () => openedMenu.value = false;

    return {
      openedMenu,

      preparedLinks,

      handleClick,
    };
  },
};
</script>
