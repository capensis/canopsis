<template>
  <c-advanced-data-table
    :headers="headers"
    :items="themes"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    :is-disabled-item="isDisabledMap"
    advanced-pagination
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
    <template #updated="{ item }">
      {{ item.updated | date }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="updatable"
          :tooltip="item.deletable ? $t('common.edit') : $t('theme.defaultTheme')"
          :disabled="!item.deletable"
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
          :tooltip="item.deletable ? $t('common.delete') : $t('theme.defaultTheme')"
          :disabled="!item.deletable"
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
    themes: {
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
          text: this.$tc('common.title'),
          value: 'name',
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
  methods: {
    isDisabledMap({ deletable }) {
      return !deletable;
    },
  },
};
</script>
