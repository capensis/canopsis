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
    template(#method="{ item }") {{ $t(`stateSetting.methods.${item.method}`) }}
    template(#actions="{ item }")
      v-layout(row)
        c-action-btn(
          v-if="updatable",
          type="edit",
          :disabled="!item.editable",
          @click.stop="$emit('edit', item)"
        )
        c-action-btn(
          v-if="addable",
          type="duplicate",
          @click.stop="$emit('duplicate', item)"
        )
        c-action-btn(
          v-if="removable",
          type="delete",
          :disabled="!item.editable",
          @click.stop="$emit('remove', item)"
        )
    template(#expand="{ item }")
      state-settings-list-expand-panel(:state-setting="item")
</template>

<script>
import StateSettingsListExpandPanel from './partials/state-settings-list-expand-panel.vue';

export default {
  components: { StateSettingsListExpandPanel },
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
          value: 'type',
          sortable: false,
        },
        {
          text: this.$t('common.enabled'),
          value: 'enabled',
          sortable: false,
        },
        {
          text: this.$t('common.priority'),
          value: 'priority',
          sortable: false,
        },
        {
          text: this.$t('common.method'),
          value: 'method',
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
