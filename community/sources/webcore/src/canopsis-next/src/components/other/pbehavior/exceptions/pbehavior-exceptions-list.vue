<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="pbehaviorExceptions",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :is-disabled-item="isDisabledException",
    select-all,
    expand,
    search,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(slot="toolbar", slot-scope="props")
      v-flex(v-show="hasDeleteAnyPbehaviorExceptionAccess && props.selected.length", xs4)
        v-btn(@click="$emit('remove-selected', props.selected)", icon)
          v-icon delete
    template(slot="actions", slot-scope="props")
      c-action-btn(
        v-if="hasUpdateAnyPbehaviorExceptionAccess",
        type="edit",
        @click="$emit('edit', props.item)"
      )
      c-action-btn(
        v-if="hasDeleteAnyPbehaviorExceptionAccess",
        :tooltip="props.item.deletable ? $t('common.delete') : $t('pbehaviorExceptions.usingException')",
        :disabled="!props.item.deletable",
        type="delete",
        @click="$emit('remove', props.item._id)"
      )
    template(slot="expand", slot-scope="props")
      pbehavior-exceptions-list-expand-panel(:pbehaviorException="props.item")
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
    pagination: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      selected: [],
    };
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
