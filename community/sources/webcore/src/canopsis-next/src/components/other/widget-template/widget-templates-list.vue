<template>
  <c-advanced-data-table
    :headers="headers"
    :items="widgetTemplates"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    advanced-pagination
    search
    @update:options="$emit('update:options', $event)"
  >
    <template #updated="{ item }">
      {{ item.updated | date }}
    </template>
    <template #type="{ item }">
      {{ $t(`widgetTemplate.types.${item.type}`) }}
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
          @click.stop="$emit('remove', item._id)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
export default {
  props: {
    widgetTemplates: {
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
      default: false,
    },
    removable: {
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
        { text: this.$t('common.name'), value: 'title', sortable: false },
        { text: this.$t('common.type'), value: 'type', sortable: false },
        { text: this.$t('common.lastModifiedOn'), value: 'updated', sortable: false },
        { text: this.$t('common.lastModifiedBy'), value: 'author.display_name', sortable: false },
        { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
      ];
    },
  },
};
</script>
