<template>
  <c-advanced-data-table
    :headers="headers"
    :items="idleRules"
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
    <template #type="{ item }">
      {{ $t(`idleRules.types.${item.type}`) }}
    </template>
    <template #operation.type="{ item }">
      {{ item | get('operation.type', '-') }}
    </template>
    <template #duration="{ item }">
      <span>{{ item.duration | duration }}</span>
    </template>
    <template #priority="{ item }">
      {{ item.priority || '-' }}
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
      </v-layout>
    </template>
    <template #expand="{ item }">
      <idle-rules-list-expand-item :idle-rule="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import IdleRulesListExpandItem from './partials/idle-rules-list-expand-item.vue';

export default {
  components: { IdleRulesListExpandItem },
  props: {
    idleRules: {
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
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.type'), value: 'type' },
        { text: this.$t('common.enabled'), value: 'enabled', sortable: false },
        { text: this.$tc('common.action'), value: 'operation.type', sortable: false },
        { text: this.$t('idleRules.timeAwaiting'), value: 'duration', sortable: false },
        { text: this.$t('common.priority'), value: 'priority' },
        { text: this.$t('common.author'), value: 'author.display_name' },
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
