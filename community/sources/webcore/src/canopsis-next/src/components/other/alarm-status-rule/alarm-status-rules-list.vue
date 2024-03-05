<template>
  <c-advanced-data-table
    :headers="headers"
    :items="rules"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    expand
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-if="removable"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #duration="{ item }">
      <span>{{ item.duration | duration }}</span>
    </template>
    <template #priority="{ item }">
      {{ item.priority || '-' }}
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
      </v-layout>
    </template>
    <template #expand="{ item }">
      <alarm-status-rules-list-expand-item :rule="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import AlarmStatusRulesListExpandItem from './partials/alarm-status-rules-list-expand-item.vue';

export default {
  components: { AlarmStatusRulesListExpandItem },
  props: {
    rules: {
      type: Array,
      required: true,
    },
    flapping: {
      type: Boolean,
      default: false,
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
        { text: this.$t('common.id'), value: '_id' },
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.duration'), value: 'duration', sortable: false },
        { text: this.$t('common.priority'), value: 'priority' },

        this.flapping && { text: this.$t('common.frequencyLimit'), value: 'freq_limit' },

        { text: this.$t('common.author'), value: 'author.display_name' },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        },
      ].filter(Boolean);
    },
  },
};
</script>
