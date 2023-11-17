<template>
  <c-advanced-data-table
    :headers="headers"
    :items="roles"
    :loading="pending"
    :options="options"
    :total-items="totalItems"
    :select-all="removable"
    advanced-pagination
    search
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-if="removable"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #auth_config.inactivity_interval="{ item }">
      {{ durationToString(item.auth_config.inactivity_interval) }}
    </template>
    <template #auth_config.expiration_interval="{ item }">
      {{ durationToString(item.auth_config.expiration_interval) }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="updatable"
          :disabled="!item.editable"
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
          :disabled="!item.deletable"
          type="delete"
          @click="$emit('remove', item._id)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
export default {
  props: {
    roles: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    options: {
      type: Object,
      required: false,
    },
    totalItems: {
      type: Number,
      required: false,
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
          text: this.$t('role.inactivityInterval'),
          value: 'auth_config.inactivity_interval',
          sortable: false,
        },
        {
          text: this.$t('role.expirationInterval'),
          value: 'auth_config.expiration_interval',
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
  methods: {
    durationToString(duration) {
      return duration ? `${duration.value}${duration.unit}` : this.$t('common.notAvailable');
    },
  },
};
</script>
