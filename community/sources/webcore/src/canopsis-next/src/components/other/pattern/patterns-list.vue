<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="patterns",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    select-all,
    search,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#mass-actions="{ selected }")
      c-action-btn(
        type="delete",
        @click="$emit('remove-selected', selected)"
      )
    template(#type="{ item }")
      span {{ $t(`pattern.types.${item.type}`) }}
    template(#updated="{ item }")
      span {{ item.updated | date }}
    template(#actions="{ item }")
      v-layout(row)
        c-action-btn(
          type="edit",
          @click="$emit('edit', item)"
        )
        c-action-btn(
          type="delete",
          @click="$emit('remove', item._id)"
        )
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
    pagination: {
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
        headers.push({ text: this.$t('common.lastModifiedBy'), value: 'author.name' });
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
