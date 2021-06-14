<template lang="pug">
  c-advanced-data-table.white(
    :headers="headers",
    :items="roles",
    :loading="pending",
    :pagination="pagination",
    :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
    :total-items="totalItems",
    advanced-pagination,
    search,
    select-all,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(slot="toolbar", slot-scope="props")
      v-flex(v-show="hasDeleteAnyRoleAccess && props.selected.length", xs4)
        v-btn(icon, data-test="massDeleteButton", @click="$emit('remove-selected', props.selected)")
          v-icon(color="error") delete
    template(slot="actions", slot-scope="props")
      v-layout(row)
        c-action-btn(
          v-if="hasUpdateAnyRoleAccess",
          type="edit",
          @click="$emit('edit', props.item)"
        )
        c-action-btn(
          v-if="hasDeleteAnyRoleAccess",
          type="delete",
          @click="$emit('remove', props.item._id)"
        )
</template>

<script>
import { permissionsTechnicalRoleMixin } from '@/mixins/permissions/technical/role';

export default {
  mixins: [permissionsTechnicalRoleMixin],
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
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.name'),
          value: 'name',
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
