<template>
  <c-advanced-data-table
    :options="options"
    :items="dynamicInfos"
    :loading="pending"
    :headers="headers"
    :total-items="totalItems"
    :select-all="removable"
    advanced-search
    advanced-pagination
    hide-actions
    expand
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, selectedKeys }">
      <c-action-btn
        v-if="removable"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
      <c-db-export-btn :ids="selectedKeys" dynamic-info />
    </template>
    <template #created="{ item }">
      {{ item.created | date }}
    </template>
    <template #updated="{ item }">
      {{ item.updated | date }}
    </template>
    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
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
          type="delete"
          @click="$emit('remove', item._id)"
        />
        <c-db-export-btn :id="item._id" dynamic-info />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <dynamic-infos-list-expand-item :info="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import DynamicInfosListExpandItem from './partials/dynamic-infos-expand-item.vue';

export default {
  components: {
    DynamicInfosListExpandItem,
  },
  props: {
    dynamicInfos: {
      type: Array,
      required: true,
    },
    options: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pending: {
      type: Boolean,
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
        { text: this.$t('common.id'), value: '_id' },
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.description'), value: 'description', sortable: false },
        { text: this.$t('common.enabled'), value: 'enabled', sortable: false },
        { text: this.$t('common.author'), value: 'author.display_name' },
        { text: this.$t('common.created'), value: 'created' },
        { text: this.$t('common.updated'), value: 'updated' },
        { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
      ];
    },
  },
};
</script>
