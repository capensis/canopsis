<template lang="pug">
  c-advanced-data-table.white(
    :headers="headers",
    :items="users",
    :loading="pending",
    :total-items="totalItems",
    :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
    :pagination="pagination",
    advanced-pagination,
    search,
    select-all,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(slot="toolbar", slot-scope="props")
      v-flex(v-show="hasDeleteAnyUserAccess && props.selected.length", xs4)
        v-btn(@click="$emit('remove-selected', props.selected)", icon)
          v-icon delete
    template(slot="enable", slot-scope="props")
      c-enabled(:value="props.item.enable")
    template(slot="source", slot-scope="props") {{ props.item.source || $constants.AUTH_METHODS.local }}
    template(slot="actions", slot-scope="props")
      v-layout(row)
        c-action-btn(
          v-if="hasUpdateAnyUserAccess",
          type="edit",
          @click.stop="$emit('edit', props.item)"
        )
        c-action-btn(
          v-if="hasDeleteAnyUserAccess",
          type="delete",
          @click.stop="$emit('remove', props.item)"
        )
</template>

<script>
import { permissionsTechnicalUserMixin } from '@/mixins/permissions/technical/user';

export default {
  mixins: [permissionsTechnicalUserMixin],
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
