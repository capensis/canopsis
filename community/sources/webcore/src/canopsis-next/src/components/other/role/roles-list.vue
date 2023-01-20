<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="roles",
    :loading="pending",
    :pagination="pagination",
    :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
    :total-items="totalItems",
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
    template(#auth_config.inactivity_interval="{ item }") {{ durationToString(item.auth_config.inactivity_interval) }}
    template(#auth_config.expiration_interval="{ item }") {{ durationToString(item.auth_config.expiration_interval) }}
    template(#actions="{ item }")
      v-layout(row)
        c-action-btn(
          v-if="updatable",
          :disabled="!item.editable",
          type="edit",
          @click="$emit('edit', item)"
        )
        c-action-btn(
          v-if="removable",
          :disabled="!item.deletable",
          type="delete",
          @click="$emit('remove', item._id)"
        )
</template>

<script>
export default {
  props: {
    roles: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    pagination: {
      type: Object,
      required: false,
    },
    totalItems: {
      type: Number,
      required: false,
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
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('role.inactivityInterval'),
          value: 'auth_config.inactivity_interval',
          sortable: false,
        },
        {
          text: this.$t('role.expirationInterval'),
          value: 'auth_config.expiration_interval',
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
  methods: {
    durationToString(duration) {
      return duration ? `${duration.value}${duration.unit}` : this.$t('common.notAvailable');
    },
  },
};
</script>
