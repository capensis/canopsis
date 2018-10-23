<template lang="pug">
  v-data-table(
  :headers="headers",
  :items="roles",
  :loading="pending"
  )
    template(
    slot="items",
    slot-scope="props",
    )
      td {{ props.item.crecord_name }}
</template>

<script>
import entitiesRoleMixins from '@/mixins/entities/role';

export default {
  mixins: [entitiesRoleMixins],
  data() {
    return {
      headers: [
        {
          text: 'Name',
          value: 'crecord_name',
        },
      ],
    };
  },
  mounted() {
    this.fetchRolesList();
  },
};
</script>
