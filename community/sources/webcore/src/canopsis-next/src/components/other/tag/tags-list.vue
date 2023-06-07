<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="tags",
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
    template(#name="{ item }")
      c-alarm-action-chip(:color="item.color") {{ item.name }}
    template(#type="{ item }")
      span {{ $t(`tag.types.${item.type}`) }}
    template(#updated="{ item }") {{ item.created | date }}
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
          :tooltip="item.deletable ? $t('common.delete') : $t('tag.importedTag')",
          :disabled="!item.deletable",
          type="delete",
          @click="$emit('remove', item._id)"
        )
</template>

<script>
export default {
  props: {
    tags: {
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
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.type'),
          value: 'type',
        },
        {
          text: this.$t('common.created'),
          value: 'created',
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
