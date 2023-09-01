<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="themes",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :select-all="removable",
    :is-disabled-item="isDisabledMap",
    advanced-pagination,
    expand,
    search,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#mass-actions="{ selected }")
      c-action-btn(
        v-show="removable",
        type="delete",
        @click="$emit('remove-selected', selected)"
      )
    template(#updated="{ item }") {{ item.updated | date }}
    template(#actions="{ item }")
      v-layout(row)
        c-action-btn(
          v-if="updatable",
          type="edit",
          @click="$emit('edit', item)"
        )
        c-action-btn(
          v-if="duplicable",
          type="duplicate",
          @click="$emit('duplicate', item)"
        )
        c-action-btn(
          v-if="removable",
          :tooltip="item.deletable ? $t('common.delete') : $t('theme.defaultTheme')",
          :disabled="!item.deletable",
          type="delete",
          @click="$emit('remove', item._id)"
        )
    template(#expand="{ item }")
      | Expand
</template>

<script>
export default {
  props: {
    themes: {
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
          text: this.$tc('common.title'),
          value: 'name',
        },
        {
          text: this.$t('common.lastModifiedOn'),
          value: 'updated',
        },
        {
          text: this.$t('common.lastModifiedBy'),
          value: 'author.display_name',
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
    isDisabledMap({ deletable }) {
      return !deletable;
    },
  },
};
</script>
