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
        small
        @click="$emit('remove-selected', selected)"
      />
    </template>

    <template #type="{ item }">
      <span>{{ $t(`pattern.types.${item.type}`) }}</span>
    </template>

    <template #updated="{ item }">
      <span>{{ dateToString(item.updated) }}</span>
    </template>

    <template #actions="{ item }">
      <c-layout>
        <c-action-btn
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
          type="delete"
          @click="$emit('remove', item._id)"
        />
      </c-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { convertDateToString } from '@/helpers/date/date';

export default {
  emits: ['update:options', 'edit', 'remove', 'remove-selected'],

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
        { title: this.$t('common.title'), key: 'title' },
        { title: this.$t('common.type'), key: 'type', sortable: false },
        { title: this.$t('common.lastModifiedOn'), key: 'updated' },
      ];

      if (this.corporate) {
        headers.push({ title: this.$t('common.lastModifiedBy'), key: 'author.display_name' });
      }

      headers.push({
        title: this.$t('common.actionsLabel'),
        key: 'actions',
        sortable: false,
      });

      return headers;
    },
  },

  methods: {
    dateToString: convertDateToString,
  },
};
</script>
