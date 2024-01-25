<template>
  <c-advanced-data-table
    :headers="headers"
    :items="remediationConfigurations"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-disabled-item="isDisabledConfiguration"
    :select-all="removable"
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-if="removable"
        type="delete"
        @click="$emit('remove-selected', selected)"
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
