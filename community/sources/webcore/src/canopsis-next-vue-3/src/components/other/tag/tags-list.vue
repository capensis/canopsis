<template>
  <c-advanced-data-table
    :headers="headers"
    :items="tags"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable && hasAnyDeletableTag"
    :is-disabled-item="isDisabledTag"
    :is-expandable-item="isCreatedTag"
    advanced-pagination
    expand
    search
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-action-btn
        v-show="removable"
        type="delete"
        small
        @click="$emit('remove-selected', selected, clearSelected)"
      />
    </template>

    <template #value="{ item }">
      <c-chip :color="item.color" class="px-2">
        {{ item.value }}
      </c-chip>
    </template>

    <template #type="{ item }">
      <span>{{ $t(`tag.types.${item.type}`) }}</span>
    </template>

    <template #created="{ item }">
      {{ formatDate(item.created) }}
    </template>

    <template #updated="{ item }">
      {{ formatDate(item.updated) }}
    </template>

    <template #actions="{ item }">
      <c-layout align="center">
        <c-action-btn v-if="updatable" type="edit" @click.stop="$emit('edit', item)" />
        <c-action-btn v-if="duplicable" type="duplicate" @click.stop="$emit('duplicate', item)" />
        <c-action-btn
          v-if="removable"
          :tooltip="item.deletable ? $t('common.delete') : $t('tag.importedTag')"
          :disabled="!item.deletable"
          type="delete"
          @click.stop="$emit('remove', item._id)"
        />
      </c-layout>
    </template>

    <template #expand="{ item }">
      <tags-list-expand-panel :tag="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { convertDateToString } from '@/helpers/date/date';
import { isCreatedTag } from '@/helpers/entities/tag/entity';

import TagsListExpandPanel from './partials/tags-list-expand-panel.vue';

export default {
  components: {
    TagsListExpandPanel,
  },

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

  emits: ['edit', 'duplicate', 'remove', 'remove-selected', 'update:options'],

  computed: {
    headers() {
      return [
        {
          text: this.$t('common.name'),
          value: 'value',
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

    hasAnyDeletableTag() {
      return this.tags.some(({ deletable }) => deletable);
    },
  },

  methods: {
    isDisabledTag({ deletable }) {
      return !deletable;
    },

    isCreatedTag(tag) {
      return isCreatedTag(tag);
    },

    formatDate(value) {
      return convertDateToString(value);
    },
  },
};
</script>
