<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="scenarios",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :select-all="removable",
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
    template(#headerCell="{ header }")
      span.pre-line.header-text {{ header.text }}
    template(#delay="{ item }")
      span {{ item.delay | duration }}
    template(#enabled="{ item }")
      c-help-icon(
        v-if="hasDeprecatedTrigger(item)",
        :text="$t('scenario.errors.deprecatedTriggerExist')",
        color="error",
        icon="error",
        top
      )
      c-enabled(v-else, :value="item.enabled")
    template(#created="{ item }") {{ item.created | date }}
    template(#updated="{ item }") {{ item.updated | date }}
    template(#actions="{ item }")
      v-layout(row)
        c-action-btn(
          v-if="updatable",
          :badge-value="isOldPattern(item)",
          :badge-tooltip="$t('pattern.oldPatternTooltip')",
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
          type="delete",
          @click="$emit('remove', item._id)"
        )
    template(#expand="{ item }")
      scenarios-list-expand-item(:scenario="item")
</template>

<script>
import { OLD_PATTERNS_FIELDS } from '@/constants';

import { isOldPattern } from '@/helpers/pattern';
import { isDeprecatedTrigger } from '@/helpers/entities/scenarios';

import { permissionsTechnicalExploitationScenarioMixin } from '@/mixins/permissions/technical/exploitation/scenario';

import ScenariosListExpandItem from './partials/scenarios-list-expand-item.vue';

export default {
  components: { ScenariosListExpandItem },
  mixins: [permissionsTechnicalExploitationScenarioMixin],
  props: {
    scenarios: {
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
          text: this.$t('common.id'),
          value: '_id',
        },
        {
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.delay'),
          value: 'delay',
          sortable: false,
        },
        {
          text: this.$t('common.priority'),
          value: 'priority',
        },
        {
          text: this.$t('common.enabled'),
          value: 'enabled',
        },
        {
          text: this.$t('common.author'),
          value: 'author.name',
        },
        {
          text: this.$t('common.created'),
          value: 'created',
        },
        {
          text: this.$t('common.updated'),
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
    hasDeprecatedTrigger(item) {
      return item.triggers.some(isDeprecatedTrigger);
    },

    isOldPattern({ actions = [] }) {
      return actions.some(action => isOldPattern(action, [OLD_PATTERNS_FIELDS.entity, OLD_PATTERNS_FIELDS.alarm]));
    },
  },
};
</script>
