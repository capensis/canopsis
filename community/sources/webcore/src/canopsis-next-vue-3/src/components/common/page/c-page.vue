<template>
  <div>
    <c-page-header :name="name">
      <slot v-if="$slots.header" name="header" />
    </c-page-header>

    <v-card class="mx-4 mb-4 mt-0">
      <slot />
    </v-card>

    <c-fab-btn
      :has-access="creatable"
      :dark="dark"
      v-bind="fabBind"
    >
      <span>{{ createTooltip }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { computed, getCurrentInstance } from 'vue';

export default {
  props: {
    name: {
      type: String,
      default: undefined,
    },
    creatable: {
      type: Boolean,
      default: false,
    },
    dark: {
      type: Boolean,
      default: false,
    },
    createTooltip: {
      type: String,
      required: false,
    },
  },

  emits: ['refresh', 'create'],

  setup(props, { emit }) {
    const instance = getCurrentInstance();
    const fabBind = computed(() => ({
      onRefresh: () => emit('refresh'),
      ...(instance?.vnode.props?.onCreate ? { onCreate: () => emit('create') } : {}),
    }));

    return {
      fabBind,
    };
  },
};
</script>
