<template lang="pug">
  c-advanced-data-table(
    :items="stateSettings",
    :headers="headers",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    advanced-pagination,
    expand,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#enabled="{ item }")
      c-enabled(:value="item.enabled")
    template(#priority="{ item }") {{ item.priority || '-' }}
    template(#method="{ item }") {{ getMethodLabel(item.method) }}
    template(#actions="{ item }")
      v-layout(row)
        c-action-btn(
          v-if="updatable",
          :disabled="!item.editable",
          type="edit",
          @click.stop="$emit('edit', item)"
        )
        c-action-btn(
          v-if="addable",
          :disabled="!isDuplicable(item)",
          type="duplicate",
          @click.stop="$emit('duplicate', item)"
        )
        c-action-btn(
          v-if="removable",
          :disabled="!item.deletable",
          type="delete",
          @click.stop="$emit('remove', item)"
        )
</template>

<script>
import { JUNIT_STATE_SETTING_ID, SERVICE_STATE_SETTING_ID } from '@/constants';

export default {
  props: {
    pagination: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    stateSettings: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    removable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.title'),
          value: 'title',
        },
        {
          text: this.$t('common.enabled'),
          value: 'enabled',
        },
        {
          text: this.$t('common.priority'),
          value: 'priority',
        },
        {
          text: this.$t('common.method'),
          value: 'method',
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
    isDuplicable(item) {
      return ![JUNIT_STATE_SETTING_ID, SERVICE_STATE_SETTING_ID].includes(item._id);
    },

    getMethodLabel(method) {
      return this.$te(`stateSetting.methods.${method}.label`)
        ? this.$t(`stateSetting.methods.${method}.label`)
        : this.$t(`stateSetting.junit.methods.${method}`);
    },
  },
};
</script>
