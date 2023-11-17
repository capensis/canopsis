<template>
  <c-advanced-data-table
    :headers="headers"
    :items="pbehaviorExceptions"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-disabled-item="isDisabledException"
    select-all
    expand
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-if="hasDeleteAnyPbehaviorExceptionAccess"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #actions="{ item: actionsItem }">
      <c-action-btn
        v-if="hasUpdateAnyPbehaviorExceptionAccess"
        type="edit"
        @click="$emit('edit', actionsItem)"
      />
      <c-action-btn
        v-if="hasDeleteAnyPbehaviorExceptionAccess"
        :tooltip="actionsItem.deletable ? $t('common.delete') : $t('pbehavior.exceptions.usingException')"
        :disabled="!actionsItem.deletable"
        type="delete"
        @click="$emit('remove', actionsItem._id)"
      />
    </template>
    <template #expand="{ item: expandItem }">
      <pbehavior-exceptions-list-expand-panel :pbehavior-exception="expandItem" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { permissionsTechnicalPbehaviorExceptionsMixin } from '@/mixins/permissions/technical/pbehavior-exceptions';

import PbehaviorExceptionsListExpandPanel from './partials/pbehavior-exceptions-list-expand-panel.vue';

export default {
  components: {
    PbehaviorExceptionsListExpandPanel,
  },
  mixins: [permissionsTechnicalPbehaviorExceptionsMixin],
  props: {
    pbehaviorExceptions: {
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
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.name'),
          value: 'name',
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
    isDisabledException({ deletable }) {
      return !deletable;
    },
  },
};
</script>
