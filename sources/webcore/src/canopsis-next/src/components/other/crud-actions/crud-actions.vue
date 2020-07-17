<template lang="pug">
  v-layout
    v-btn(
      v-for="action in actions",
      :key="action",
      icon,
      small,
      @click.stop="$emit(action)"
    )
      v-icon(v-bind="actionProps[action]") {{ actionProps[action].icon }}
</template>

<script>
import { CRUD_ACTIONS } from '@/constants';

export default {
  props: {
    actions: {
      type: Array,
      default: () => [CRUD_ACTIONS.update, CRUD_ACTIONS.delete],
    },
  },
  computed: {
    actionProps() {
      return {
        [CRUD_ACTIONS.create]: {
          icon: 'create',
        },
        [CRUD_ACTIONS.update]: {
          icon: 'create',
        },
        [CRUD_ACTIONS.delete]: {
          class: 'error--text',
          icon: 'delete',
        },
      };
    },
  },
};
</script>
