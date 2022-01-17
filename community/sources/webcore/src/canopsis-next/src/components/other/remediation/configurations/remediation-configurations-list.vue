<template lang="pug">
  c-advanced-data-table.white(
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
      c-action-btn(
        v-if="hasUpdateAnyRemediationConfigurationAccess",
        type="edit",
        @click="$emit('edit', props.item)"
      )
      c-action-btn(
        v-if="hasDeleteAnyRemediationConfigurationAccess",
        :tooltip="props.disabled ? $t('remediationConfigurations.usingConfiguration') : $t('common.delete')",
        :disabled="props.disabled",
        type="delete",
        @click="$emit('remove', props.item)"
      )
</template>

<script>
import {
  permissionsTechnicalRemediationConfigurationMixin,
} from '@/mixins/permissions/technical/remediation-configuration';

export default {
  mixins: [permissionsTechnicalRemediationConfigurationMixin],
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
          value: 'author.name',
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
