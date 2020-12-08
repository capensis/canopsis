<template lang="pug">
  advanced-data-table.white(
    :headers="headers",
    :items="users",
    :loading="pending",
    :total-items="totalItems",
    :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
    :pagination="pagination",
    search,
    select-all,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(slot="toolbar", slot-scope="props")
      v-flex(v-show="hasDeleteAnyUserAccess && props.selected.length", xs4)
        v-btn(@click="$emit('remove-selected', props.selected)", icon)
          v-icon delete
    template(slot="enable", slot-scope="props")
      enabled-column(:value="props.item.enable")
    template(slot="actions", slot-scope="props")
      div
        v-btn.mx-0(
          v-if="hasUpdateAnyUserAccess",
          icon,
          @click.stop="$emit('edit', props.item)"
        )
          v-icon edit
        v-btn.mx-0(
          v-if="hasDeleteAnyUserAccess",
          icon,
          @click.stop="$emit('remove', props.item)"
        )
          v-icon(color="error") delete
</template>

<script>
import rightsTechnicalUserMixin from '@/mixins/rights/technical/user';

import EnabledColumn from '@/components/tables/enabled-column.vue';

export default {
  components: {
    EnabledColumn,
  },
  mixins: [rightsTechnicalUserMixin],
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
          text: this.$t('tables.admin.users.columns.username'),
          value: 'id',
        },
        {
          text: this.$t('tables.admin.users.columns.firstName'),
          value: 'firstname',
        },
        {
          text: this.$t('tables.admin.users.columns.lastName'),
          value: 'lastname',
        },
        {
          text: this.$t('tables.admin.users.columns.role'),
          value: 'role',
        },
        {
          text: this.$t('tables.admin.users.columns.enabled'),
          value: 'enable',
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
