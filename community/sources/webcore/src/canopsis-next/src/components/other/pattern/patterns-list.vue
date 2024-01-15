<template>
  <c-advanced-data-table
    :headers="headers"
    :items="patterns"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    select-all
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #type="{ item }">
      <span>{{ $t(`pattern.types.${item.type}`) }}</span>
    </template>
    <template #updated="{ item }">
      <span>{{ item.updated | date }}</span>
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
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
    patterns: {
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
    corporate: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    headers() {
      const headers = [
        { text: this.$t('common.title'), value: 'title' },
        { text: this.$t('common.type'), value: 'type', sortable: false },
        { text: this.$t('common.lastModifiedOn'), value: 'updated' },
      ];

      if (this.corporate) {
        headers.push({ text: this.$t('common.lastModifiedBy'), value: 'author.display_name' });
      }

      headers.push({
        text: this.$t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      });

      return headers;
    },
  },
};
</script>
