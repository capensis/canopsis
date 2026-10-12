<template>
  <c-advanced-data-table
    :items="icons"
    :headers="headers"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #icon="{ item }">
      <custom-icon :content="item.content" />
    </template>
    <template #updated="{ item }">
      {{ item.updated | date }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="updatable"
          type="edit"
          @click.stop="$emit('edit', item)"
        />
        <c-action-btn
          v-if="removable"
          type="delete"
          @click.stop="$emit('remove', item)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import CustomIcon from './partials/custom-icon.vue';

export default {
  components: { CustomIcon },
  props: {
    options: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    icons: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    removable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    headers() {
      return [
        {
          text: this.$tc('common.icon', 1),
          value: 'icon',
          sortable: false,
        },
        {
          text: this.$t('common.title'),
          value: 'title',
        },
        {
          text: this.$t('common.lastModifiedOn'),
          value: 'updated',
        },
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
