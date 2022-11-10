<template lang="pug">
  div.instruction-list
    c-advanced-data-table(
      :headers="headers",
      :items="remediationInstructions",
      :loading="pending",
      :total-items="totalItems",
      :pagination="pagination",
      :is-disabled-item="isDisabledInstruction",
      :select-all="removable",
      search,
      advanced-pagination,
      @update:pagination="$emit('update:pagination', $event)"
    )
      template(#mass-actions="{ selected }")
        c-action-btn.ml-3(v-if="removable", type="delete", @click="$emit('remove-selected', selected)")
      template(#headerCell="{ header }")
        span.c-table-header__text--multiline {{ header.text }}
      template(#author="{ item }")
        span {{ item.author.name }}
      template(#enabled="{ item }")
        c-enabled(:value="item.enabled")
      template(#status="{ item }")
        v-tooltip(v-if="item.approval", bottom)
          template(#activator="{ on }")
            v-icon(color="black") query_builder
          span {{ $t('remediationInstructions.approvalPending') }}
        v-icon(v-else, color="primary") check_circle
      template(#type="{ item }") {{ $t(`remediationInstructions.types.${item.type}`) }}
      template(#last_modified="{ item }") {{ item.last_modified | date }}
      template(#last_executed_on="{ item }") {{ item.last_executed_on | date }}
      template(#actions="{ item, disabled }")
        v-layout(row, justify-end)
          c-action-btn(
            v-if="item.approval && isApprovalForCurrentUser(item.approval)",
            :tooltip="$t('remediationInstructions.needApprove')",
            icon="notification_important",
            color="error",
            @click="$emit('approve', item)"
          )
          c-action-btn(
            v-if="updatable",
            type="edit",
            @click="$emit('edit', item)"
          )
          c-action-btn(
            v-if="updatable",
            :tooltip="$t('modals.patterns.title')",
            icon="assignment",
            @click="$emit('assign-patterns', item)"
          )
          c-action-btn(
            v-if="duplicable",
            type="duplicate",
            @click="$emit('duplicate', item)"
          )
          c-action-btn(
            v-if="removable",
            :tooltip="disabled ? $t('remediationInstructions.usingInstruction') : $t('common.delete')",
            :disabled="disabled",
            type="delete",
            @click="$emit('remove', item)"
          )
</template>

<script>
import { authMixin } from '@/mixins/auth';

export default {
  mixins: [authMixin],
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
      return remediationInstruction?.user?._id === this.currentUser._id
        || remediationInstruction?.role?._id === this.currentUser.role._id;
    },

    isDisabledInstruction({ deletable }) {
      return !deletable;
    },
  },
};
</script>
