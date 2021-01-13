<template lang="pug">
  advanced-data-table.white(
    :headers="headers",
    :items="remediationConfigurations",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :is-disabled-item="isDisabledConfiguration",
    select-all,
    search,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(slot="toolbar", slot-scope="props")
      v-flex(v-show="hasDeleteAnyRemediationConfigurationAccess && props.selected.length", xs4)
        v-btn(@click="$emit('remove-selected', props.selected)", icon)
          v-icon delete
    template(slot="actions", slot-scope="props")
      action-btn(
        v-if="hasUpdateAnyRemediationConfigurationAccess",
        type="edit",
        @click="$emit('edit', props.item)"
      )
      action-btn(
        v-if="hasDeleteAnyRemediationConfigurationAccess",
        :tooltip="props.disabled ? $t('remediationConfigurations.usingConfiguration') : $t('common.delete')",
        :disabled="props.disabled",
        type="delete",
        @click="$emit('edit', props.item)"
      )
</template>

<script>
import rightsTechnicalRemediationConfigurationMixin from '@/mixins/rights/technical/remediation-configuration';

import ActionBtn from '@/components/common/buttons/action-btn.vue';

export default {
  components: { ActionBtn },
  mixins: [rightsTechnicalRemediationConfigurationMixin],
  props: {
    remediationConfigurations: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
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
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.author'),
          value: 'author',
        },
        {
          text: this.$t('common.type'),
          value: 'type',
        },
        {
          text: this.$t('remediationConfigurations.table.host'),
          value: 'host',
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
    isDisabledConfiguration({ deletable }) {
      return !deletable;
    },
  },
};
</script>
