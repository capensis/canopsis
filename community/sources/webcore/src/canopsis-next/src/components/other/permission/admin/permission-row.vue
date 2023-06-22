<template lang="pug">
  tr
    td
      v-layout(align-center)
        span.mr-1 {{ permission.name }}
        c-help-icon(
          v-if="permission.description",
          :text="permission.description",
          icon="help",
          size="18",
          top
        )
    permission-row-cell(
      v-for="role in roles",
      :key="`role-permission-${role._id}`",
      :permission="permission",
      :role="role",
      :changed-role="changedRoles[role._id]",
      :disabled="disabled",
      @change="$listeners.change"
    )
</template>

<script>
import PermissionRowCell from './permission-row-cell.vue';

export default {
  components: { PermissionRowCell },
  props: {
    permission: {
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
};
</script>
