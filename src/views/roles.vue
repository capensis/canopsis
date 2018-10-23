<template lang="pug">
  div Roles
</template>

<script>
import entitiesRoleMixins from '@/mixins/entities/role';

export default {
  mixins: [entitiesRoleMixins],
  mounted() {
    this.fetchRolesList();
  },
};
</script>
