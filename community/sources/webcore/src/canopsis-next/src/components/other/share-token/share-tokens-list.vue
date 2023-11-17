<template>
  <c-advanced-data-table
    :items="shareTokens"
    :headers="headers"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    advanced-pagination
    search
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-if="removable"
        :tooltip="$t('shareToken.revokeSelectedTokens')"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #created="{ item }">
      {{ item.created | date }}
    </template>
    <template #expired="{ item }">
      {{ item.expired | date('long', $t('common.notAvailable')) }}
    </template>
    <template #accessed="{ item }">
      {{ item.accessed | date('long', $t('common.notAvailable')) }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="removable"
          :tooltip="$t('shareToken.revokeToken')"
          type="delete"
          @click="$emit('remove', item)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
export default {
  props: {
    shareTokens: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: true,
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
  },
  computed: {
    headers() {
      const headers = [
        { text: this.$t('common.id'), value: '_id', sortable: false },
        { text: this.$t('common.description'), value: 'description' },
        { text: this.$t('common.value'), value: 'value', sortable: false },
        { text: this.$tc('common.user'), value: 'user.display_name' },
        { text: this.$tc('common.role'), value: 'role.name', sortable: false },
        { text: this.$t('common.created'), value: 'created' },
        { text: this.$t('common.expired'), value: 'expired' },
        { text: this.$t('common.accessed'), value: 'accessed' },
      ];

      if (this.removable) {
        headers.push({ text: this.$t('common.actionsLabel'), value: 'actions', sortable: false });
      }

      return headers;
    },
  },
};
</script>
