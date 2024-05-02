<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="linkRules",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :select-all="removable",
    search,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#mass-actions="{ selected }")
      c-action-btn.ml-3(v-if="removable", type="delete", @click="$emit('remove-selected', selected)")
    template(#enabled="{ item }")
      c-enabled(:value="item.enabled")
    template(#created="{ item }") {{ item.created | date }}
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
          type="delete",
          @click="$emit('remove', item._id)"
        )
</template>

<script>
export default {
  props: {
    linkRules: {
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
    duplicable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    headers() {
      return [
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.enabled'), value: 'enabled', sortable: false },
        { text: this.$t('common.lastModifiedOn'), value: 'updated' },
        { text: this.$t('common.lastModifiedBy'), value: 'author.display_name' },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        },
      ];
    },
  },
};
</script>
