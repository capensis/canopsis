<template lang="pug">
  c-advanced-data-table(
    :items="stateSettings",
    :headers="headers",
    :loading="pending",
    :total-items="stateSettings.length",
    no-pagination,
    expand
  )
    template(slot="method", slot-scope="props") {{ $t(`stateSetting.methods.${props.item.method}`) }}
    template(slot="actions", slot-scope="props")
      v-layout(row)
        c-action-btn(
          type="edit",
          :disabled="!props.item.editable",
          @click.stop="$emit('edit', props.item)"
        )
    template(slot="expand", slot-scope="props")
      state-settings-list-expand-panel(:state-setting="props.item")
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
