<template>
  <c-advanced-data-table
    :items="stateSettings"
    :headers="headers"
    :loading="pending"
    :total-items="stateSettings.length"
    no-pagination
    expand
  >
    <template #method="{ item }">
      {{ $t(`stateSetting.methods.${item.method}`) }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          type="edit"
          :disabled="!item.editable"
          @click.stop="$emit('edit', item)"
        />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <state-settings-list-expand-panel :state-setting="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import StateSettingsListExpandPanel from './partials/state-settings-list-expand-panel.vue';

export default {
  components: { StateSettingsListExpandPanel },
  props: {
    stateSettings: {
      type: Array,
      default: () => [],
    },
    pending: {
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
