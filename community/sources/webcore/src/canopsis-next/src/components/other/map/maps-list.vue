<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="maps",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :select-all="removable",
    :is-disabled-item="isDisabledMap",
    advanced-pagination,
    expand,
    search,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#mass-actions="{ selected }")
      c-action-btn(
        v-show="removable",
        type="delete",
        @click="$emit('remove-selected', selected)"
      )
    template(#type="{ item }")
      span {{ $t(`map.types.${item.type}`) }}
    template(#updated="{ item }") {{ item.updated | date }}
    template(#actions="{ item }")
      v-layout(row)
        c-action-btn(
          v-if="updatable",
          type="edit",
          @click="$emit('edit', item)"
        )
        c-action-btn(
          v-if="duplicable",
          type="duplicate",
          @click="$emit('duplicate', item)"
        )
        c-action-btn(
          v-if="removable",
          :tooltip="item.deletable ? $t('common.delete') : $t('map.usingMap')",
          :disabled="!item.deletable",
          type="delete",
          @click="$emit('remove', item._id)"
        )
    template(#expand="{ item }")
      maps-list-expand-item(:map="item")
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
    pagination: {
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
          value: 'author.name',
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
