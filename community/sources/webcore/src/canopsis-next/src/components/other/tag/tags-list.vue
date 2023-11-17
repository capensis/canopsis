<template>
  <c-advanced-data-table
    :headers="headers"
    :items="tags"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    :is-disabled-item="isDisabledTag"
    :is-expandable-item="isCreatedTag"
    advanced-pagination
    expand
    search
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-show="removable"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #value="{ item }">
      <c-alarm-action-chip :color="item.color">
        {{ item.value }}
      </c-alarm-action-chip>
    </template>
    <template #type="{ item }">
      <span>{{ $t(`tag.types.${item.type}`) }}</span>
    </template>
    <template #created="{ item }">
      {{ item.created | date }}
    </template>
    <template #updated="{ item }">
      {{ item.updated | date }}
    </template>
    <template #actions="{ item }">
      <v-layout>
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
          :tooltip="item.deletable ? $t('common.delete') : $t('tag.importedTag')"
          :disabled="!item.deletable"
          type="delete"
          @click="$emit('remove', item._id)"
        />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <tags-list-expand-panel :tag="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { isCreatedTag } from '@/helpers/entities/tag/entity';

import TagsListExpandPanel from './partials/tags-list-expand-panel.vue';

export default {
  components: { TagsListExpandPanel },
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
  },
  methods: {
    isDisabledTag({ deletable }) {
      return !deletable;
    },

    isCreatedTag(tag) {
      return isCreatedTag(tag);
    },
  },
};
</script>
