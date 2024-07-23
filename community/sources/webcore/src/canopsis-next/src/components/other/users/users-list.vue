<template lang="pug">
  c-advanced-data-table(
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
      c-action-btn(
        v-if="removable",
        type="delete",
        @click="$emit('remove-selected', selected)"
      )
    template(#enable="{ item }")
      c-enabled(:value="item.enable")
    template(#active="{ item }")
      c-enabled(:value="item.active_connects > 0")
    template(#source="{ item }") {{ item.source || $constants.AUTH_METHODS.local }}
    template(#roles="{ item }")
      v-chip-group(:items="item.roles", item-text="name", item-value="_id")
    template(#actions="{ item }")
      v-layout(row)
        c-action-btn(
          v-if="updatable",
          type="edit",
          @click.stop="$emit('edit', item)"
        )
        c-action-btn(
          v-if="removable",
          :disabled="!item.deletable",
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
          text: this.$t('common.username'),
          value: 'name',
        },
        {
          text: this.$t('user.displayName'),
          value: 'display_name',
        },
        {
          text: this.$t('user.firstName'),
          value: 'firstname',
          sortable: false,
        },
        {
          text: this.$t('user.lastName'),
          value: 'lastname',
          sortable: false,
        },
        {
          text: this.$tc('common.role', 2),
          value: 'roles',
          sortable: false,
        },
        {
          text: this.$t('user.active'),
          value: 'active',
          sortable: false,
        },
        {
          text: this.$t('common.enabled'),
          value: 'enable',
        },
        {
          text: this.$t('user.auth'),
          value: 'source',
        },
        {
          text: this.$t('user.activeConnects'),
          value: 'active_connects',
          sortable: false,
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
