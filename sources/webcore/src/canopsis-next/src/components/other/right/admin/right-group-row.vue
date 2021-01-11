<template lang="pug">
  tr
    td.cursor-pointer
      expand-button.mr-2(:expanded="expanded", @expand="$emit('expand')")
      span {{ group.name }}
    right-group-row-cell(
      v-for="role in roles",
      :key="`role-right-${role._id}`",
      :group="group",
      :role="role",
      :changedRole="changedRoles[role._id]",
      :disabled="disabled",
      @change="change"
    )
</template>

<script>
import ExpandButton from '@/components/common/buttons/expand-button.vue';

import RightGroupRowCell from './right-group-row-cell.vue';

export default {
  components: { ExpandButton, RightGroupRowCell },
  props: {
    group: {
      type: Object,
      required: true,
    },
    roles: {
      type: Array,
      default: () => [],
    },
    changedRoles: {
      type: Object,
      default: () => ({}),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    expanded: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    change(...args) {
      this.$emit('change', ...args);
    },
  },
};
</script>
