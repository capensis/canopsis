<template>
  <c-advanced-data-table
    :headers="headers"
    :items="remediationConfigurations"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-disabled-item="isDisabledConfiguration"
    :select-all="removable"
    :hide-mass-actions="!active"
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        small
        remediation-configuration
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>
    <template #actions="{ item, disabled }">
      <c-action-btn
        v-if="updatable"
        type="edit"
        @click="$emit('edit', item)"
      />
      <c-action-btn
        v-if="duplicable"
        type="duplicate"
        @click="$emit('duplicate', item)"
      />
      <c-action-btn
        v-if="removable"
        :tooltip="disabled ? $t('remediation.configuration.usingConfiguration') : $t('common.delete')"
        :disabled="disabled"
        type="delete"
        @click="$emit('remove', item)"
      />
      <c-db-export-btn :id="item._id" job-config />
    </template>
  </c-advanced-data-table>
</template>

<script>
export default {
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
    options: {
      type: Object,
      required: true,
    },
    removable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    duplicable: {
      type: Boolean,
      default: false,
    },
    active: {
      type: Boolean,
      default: true,
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
          value: 'author.display_name',
        },
        {
          text: this.$t('common.type'),
          value: 'type',
        },
        {
          text: this.$t('remediation.configuration.host'),
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
