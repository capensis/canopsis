<template>
  <c-advanced-data-table
    :headers="headers"
    :items="scenarios"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    expand
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-if="removable"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #headerCell="{ header }">
      <span class="pre-line header-text">{{ header.text }}</span>
    </template>
    <template #delay="{ item }">
      <span>{{ item.delay | duration }}</span>
    </template>
    <template #priority="{ item }">
      {{ item.priority || '-' }}
    </template>
    <template #enabled="{ item }">
      <c-help-icon
        v-if="hasDeprecatedTrigger(item)"
        :text="$t('scenario.errors.deprecatedTriggerExist')"
        color="error"
        icon="error"
        top
      />
      <c-enabled
        v-else
        :value="item.enabled"
      />
    </template>
    <template #created="{ item }">
      {{ item.created | date }}
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
          type="delete"
          @click="$emit('remove', item._id)"
        />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <scenarios-list-expand-item :scenario="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { isDeprecatedTrigger } from '@/helpers/entities/scenario/form';

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
          value: 'author.display_name',
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
      return item.triggers.some(({ type }) => isDeprecatedTrigger(type));
    },
  },
};
</script>
