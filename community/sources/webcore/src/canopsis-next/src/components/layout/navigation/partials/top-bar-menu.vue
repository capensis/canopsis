<template>
  <v-menu
    v-if="preparedLinks.length"
    v-field="value"
    v-bind="$attrs"
    :position-x="positionX"
    :position-y="positionY"
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
    />
  </v-menu>
</template>

<script>
import { computed, toRef } from 'vue';

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
      default: true,
    },
    positionX: {
      type: Number,
      default: 0,
    },
    positionY: {
      type: Number,
      default: 0,
    },
  },
  setup(props) {
    const { prepareLinks } = useTopBarMenu({
      withoutSort: toRef(props, 'withoutSort'),
      permissionsWithDefaultType: toRef(props, 'permissionsWithDefaultType'),
    });

    const preparedLinks = computed(() => prepareLinks(props.links));

    return {
      preparedLinks,
    };
  },
};
</script>
