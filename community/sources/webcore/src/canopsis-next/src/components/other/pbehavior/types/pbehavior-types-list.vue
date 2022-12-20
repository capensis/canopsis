<template lang="pug">
  c-advanced-data-table.white(
    :headers="headers",
    :items="pbehaviorTypes",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :is-disabled-item="isDisabledType",
    :select-all="deletable",
    expand,
    search,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#mass-actions="{ selected }")
      c-action-btn(
        v-if="removable",
        type="delete",
        @click="$emit('remove-selected', selected)"
      )
    template(#icon_name="{ item }")
      span.pbehavior-type-icon.secondary(v-if="item.icon_name")
        v-icon(color="white", size="18") {{ item.icon_name }}
    template(#actions="{ item }")
      v-layout
        c-action-btn(
          :disabled="!item.editable",
          :tooltip="item.editable ? $t('common.edit') : $t('pbehaviorTypes.defaultType')",
          type="edit",
          @click="$emit('edit', item)"
        )
        c-action-btn(
          :disabled="!item.deletable",
          :tooltip="item.deletable ? $t('common.delete') : $t('pbehaviorTypes.defaultType')",
          type="delete",
          @click="$emit('remove', item._id)"
        )
    template(#expand="{ item }")
      pbehavior-types-list-expand-panel(:pbehavior-type="item")
</template>

<script>
import PbehaviorTypesListExpandPanel from './partials/pbehavior-types-list-expand-panel.vue';

export default {
  components: {
    PbehaviorTypesListExpandPanel,
  },
  props: {
    pbehaviorTypes: {
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
      required: true,
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
          text: this.$t('common.icon'),
          value: 'icon_name',
          sortable: false,
        },
        {
          text: this.$t('common.priority'),
          value: 'priority',
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
    isDisabledType({ deletable }) {
      return !deletable;
    },
  },
};
</script>

<style lang="scss" scoped>
  .pbehavior-type-icon {
    display: inline-flex;
    padding: 2px 10px;
    border-radius: 10px;
  }
</style>
