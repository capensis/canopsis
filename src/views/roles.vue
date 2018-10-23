<template lang="pug">
  v-container
    div
      div(v-for="role in roles", :key="role._id")
        strong {{ role.crecord_name }}
    v-btn(@click="showCreateRoleModal") Add role
</template>

<script>
import entitiesRoleMixins from '@/mixins/entities/role';
import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

export default {
  mixins: [entitiesRoleMixins, modalMixin],
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
  methods: {
    showCreateRoleModal() {
      this.showModal({
        name: MODALS.createRole,
        config: {
          title: 'modals.createRole.title',
        },
      });
    },
  },
};
</script>
