<template lang="pug">
  tr
    td {{ right.desc }}
    right-row-cell(
      v-for="role in roles",
      :key="`role-right-${role._id}`",
      :right="right",
      :role="role",
      :changedRole="changedRoles[role._id]",
      :disabled="disabled",
      @change="change"
    )
</template>

<script>
import RightRowCell from './right-row-cell.vue';

export default {
  components: { RightRowCell },
  props: {
    right: {
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
  },
  methods: {
    change(...args) {
      this.$emit('change', ...args);
    },
  },
};
</script>
