<template lang="pug">
  c-advanced-data-table.white(
    :headers="headers",
    :items="users",
    :loading="pending",
    :total-items="totalItems",
    :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
    :pagination="pagination",
    :select-all="removable",
    advanced-pagination,
    search,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#mass-actions="{ selected }")
      c-action-btn.ml-3(v-if="removable", type="delete", @click="$emit('remove-selected', selected)")
    template(#enable="{ item }")
      c-enabled(:value="item.enable")
    template(#source="{ item }") {{ item.source || $constants.AUTH_METHODS.local }}
    template(#actions="{ item }")
      v-layout(row)
        c-action-btn(
          v-if="updatable",
          type="edit",
          @click.stop="$emit('edit', item)"
        )
        c-action-btn(
          v-if="removable",
          type="delete",
          @click.stop="$emit('remove', item)"
        )
</template>

<script>
export default {
  props: {
    users: {
      type: Array,
      default: () => [],
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    pagination: {
      type: Object,
      required: true,
    },
    removable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('users.username'),
          value: 'name',
        },
        {
          text: this.$t('users.firstName'),
          value: 'firstname',
          sortable: false,
        },
        {
          text: this.$t('users.lastName'),
          value: 'lastname',
          sortable: false,
        },
        {
          text: this.$t('users.role'),
          value: 'role.name',
        },
        {
          text: this.$t('users.enabled'),
          value: 'enable',
        },
        {
          text: this.$t('users.auth'),
          value: 'source',
        },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        },
      ];
    },
  },
};
</script>
