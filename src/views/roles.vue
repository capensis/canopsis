<template lang="pug">
  v-container
    v-layout(justify-end)
      v-btn(@click="showCreateRoleModal", icon, fab)
        v-icon(color="green darken-4") add
    v-data-table(
    :items="roles",
    :headers="headers"
    )
      template(slot="headers", slot-scope="props")
        tr
          th(v-for="header in headers", :key="header.value") {{ header.text }}
      template(slot="items", slot-scope="props")
        td.text-xs-center {{ props.item.crecord_name }}
        td.text-xs-center
          v-btn(icon)
            v-icon edit
          v-btn(icon)
            v-icon(color="red darken-4") delete
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
        {
          text: 'Actions',
          value: 'actions',
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
