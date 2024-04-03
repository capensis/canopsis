<template>
  <c-advanced-data-table
    :headers="headers"
    :items="maps"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    :is-disabled-item="isDisabledMap"
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
    <template #type="{ item }">
      <span>{{ $t(`map.types.${item.type}`) }}</span>
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
          :tooltip="item.deletable ? $t('common.delete') : $t('map.usingMap')"
          :disabled="!item.deletable"
          type="delete"
          @click="$emit('remove', item._id)"
        />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <maps-list-expand-item :map="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import MapsListExpandItem from './partials/maps-list-expand-item.vue';

export default {
  components: {
    MapsListExpandItem,
  },
  props: {
    maps: {
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
          value: 'name',
        },
        {
          text: this.$t('common.type'),
          value: 'type',
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
    isDisabledMap({ deletable }) {
      return !deletable;
    },
  },
};
</script>
