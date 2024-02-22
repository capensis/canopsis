<template>
  <tr>
    <td class="cursor-pointer">
      <c-expand-btn
        :expanded="expanded"
        class="mr-2"
        @expand="$emit('expand', !expanded)"
      />
      <span>{{ group.name }}</span>
    </td>
    <permission-group-row-cell
      v-for="role in roles"
      :key="`role-permission-${role._id}`"
      :group="group"
      :role="role"
      :changed-role="changedRoles[role._id]"
      :disabled="disabled || isEmptyPermissions"
      @change="change"
    />
  </tr>
</template>

<script>
import PermissionGroupRowCell from './permission-group-row-cell.vue';

export default {
  components: { PermissionGroupRowCell },
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
  computed: {
    isEmptyPermissions() {
      return !this.group?.permissions?.length;
    },
  },
  methods: {
    change(...args) {
      this.$emit('change', ...args);
    },
  },
};
</script>
