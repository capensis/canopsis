<template lang="pug">
  tr
    td.cursor-pointer
      c-expand-btn.mr-2(:expanded="expanded", @expand="$emit('expand')")
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
import RightGroupRowCell from './right-group-row-cell.vue';

export default {
  components: { RightGroupRowCell },
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
