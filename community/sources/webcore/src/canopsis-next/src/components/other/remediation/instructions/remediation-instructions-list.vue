<template lang="pug">
  div.instruction-list
    c-advanced-data-table.white(
      :headers="headers",
      :items="remediationInstructions",
      :loading="pending",
      :total-items="totalItems",
      :pagination="pagination",
      :is-disabled-item="isDisabledInstruction",
      select-all,
      search,
      advanced-pagination,
      @update:pagination="$emit('update:pagination', $event)"
    )
      template(slot="toolbar", slot-scope="props")
        v-flex(v-show="hasDeleteAnyRemediationInstructionAccess && props.selected.length", xs4)
          v-btn(@click="$emit('remove-selected', props.selected)", icon)
            v-icon delete
      template(slot="headerCell", slot-scope="props")
        span.c-table-header__text--multiline {{ props.header.text }}
      template(slot="author", slot-scope="props")
        span {{ props.item.author.name }}
      template(slot="enabled", slot-scope="props")
        c-enabled(:value="props.item.enabled")
      template(slot="status", slot-scope="props")
        v-tooltip(v-if="props.item.approval", bottom)
          slot(slot="activator")
            v-icon(color="black") query_builder
          span {{ $t('remediationInstructions.approvalPending') }}
        v-icon(v-else, color="primary") check_circle
      template(slot="type", slot-scope="props") {{ $t(`remediationInstructions.types.${props.item.type}`) }}
      template(slot="last_modified", slot-scope="props") {{ props.item.last_modified | date }}
      template(slot="last_executed_on", slot-scope="props") {{ props.item.last_executed_on | date }}
      template(slot="actions", slot-scope="props")
        v-layout(row, justify-end)
          c-action-btn(
            v-if="props.item.approval && isApprovalForCurrentUser(props.item.approval)",
            :tooltip="$t('remediationInstructions.needApprove')",
            icon="notification_important",
            color="error",
            @click="$emit('approve', props.item)"
          )
          c-action-btn(
            v-if="hasUpdateAnyRemediationInstructionAccess",
            type="edit",
            @click="$emit('edit', props.item)"
          )
          c-action-btn(
            v-if="hasUpdateAnyRemediationInstructionAccess",
            :tooltip="$t('modals.patterns.title')",
            icon="assignment",
            @click="$emit('assign-patterns', props.item)"
          )
          c-action-btn(
            v-if="hasDeleteAnyRemediationInstructionAccess",
            :tooltip="props.disabled ? $t('remediationInstructions.usingInstruction') : $t('common.delete')",
            :disabled="props.disabled",
            type="delete",
            @click="$emit('remove', props.item)"
          )
</template>

<script>
import { get } from 'lodash';

import {
  permissionsTechnicalRemediationInstructionMixin,
} from '@/mixins/permissions/technical/remediation-instruction';

export default {
  mixins: [permissionsTechnicalRemediationInstructionMixin],
  props: {
    remediationInstructions: {
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
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.author'),
          value: 'author',
        },
        {
          text: this.$t('common.enabled'),
          value: 'enabled',
        },
        {
          text: this.$t('common.type'),
          value: 'type',
        },
        {
          text: this.$t('common.lastModifiedOn'),
          value: 'last_modified',
        },
        {
          text: this.$t('common.status'),
          value: 'status',
        },
        {
          text: this.$t('remediationInstructions.table.monthExecutions'),
          value: 'month_executions',
          sortable: false,
        },
        {
          text: this.$t('remediationInstructions.table.lastExecutedOn'),
          value: 'last_executed_on',
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
    isApprovalForCurrentUser(remediationInstruction) {
      return get(remediationInstruction, 'user._id') === this.currentUser._id
        || get(remediationInstruction, 'role._id') === this.currentUser.role._id;
    },

    isDisabledInstruction({ deletable }) {
      return !deletable;
    },
  },
};
</script>
