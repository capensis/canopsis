<template>
  <c-advanced-data-table
    :headers="headers"
    :items="linkRules"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
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
    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
    </template>
    <template #author.name="{ item }">
      {{ item.author?.display_name }}
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
          type="delete"
          @click="$emit('remove', item._id)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
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
    options: {
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
        { text: this.$t('common.lastModifiedBy'), value: 'author.name' },
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
