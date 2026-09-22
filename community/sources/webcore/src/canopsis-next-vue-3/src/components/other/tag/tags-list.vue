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
      <c-chip :color="item.color" :text-color="getTagTextColor(item.color)" class="px-2">
        {{ item.value }}
      </c-chip>
    </template>

    <template #type="{ item }">
      <span>{{ t(`tag.types.${item.type}`) }}</span>
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
          :tooltip="item.deletable ? t('common.delete') : t('tag.importedTag')"
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
import { computed } from 'vue';

import { getMostReadableTextColor } from '@/helpers/color';
import { convertDateToString } from '@/helpers/date/date';
import { isCreatedTag } from '@/helpers/entities/tag/entity';

import { useI18n } from '@/hooks/i18n';

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

  setup(props) {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.name'),
        value: 'value',
        sortable: true,
      },
      {
        text: t('common.type'),
        value: 'type',
        sortable: false,
      },
      {
        text: t('common.created'),
        value: 'created',
        sortable: true,
      },
      {
        text: t('common.lastModifiedOn'),
        value: 'updated',
        sortable: false,
      },
      {
        text: t('common.lastModifiedBy'),
        value: 'author.display_name',
        sortable: false,
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    const hasAnyDeletableTag = computed(() => props.tags.some(({ deletable }) => deletable));

    const isDisabledTag = ({ deletable }) => !deletable;
    const getTagTextColor = color => getMostReadableTextColor(color, { level: 'AA', size: 'large' });
    const formatDate = value => convertDateToString(value);

    return {
      t,
      headers,
      hasAnyDeletableTag,
      isDisabledTag,
      isCreatedTag,
      getTagTextColor,
      formatDate,
    };
  },
};
</script>
