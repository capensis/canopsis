<template>
  <v-layout column>
    <view-duplicate-private-field
      v-if="duplicatePrivate"
      v-field="form.is_private"
    />
    <c-name-field
      v-field="form.title"
      :label="$t('common.title')"
      name="title"
      required
    />
    <v-text-field
      v-field="form.description"
      :label="$t('common.description')"
      name="description"
    />
    <c-enabled-field v-field="form.enabled" />
    <periodic-refresh-field
      v-field="form.periodic_refresh"
      :label="$t('view.periodicRefresh')"
    />
    <view-tags-field v-field="form.tags" />
    <view-group-field
      v-field="form.group"
      :groups="groups"
    />
  </v-layout>
</template>

<script>
import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';

import ViewTagsField from './fields/view-tags-field.vue';
import ViewGroupField from './fields/view-group-field.vue';
import ViewDuplicatePrivateField from './fields/view-duplicate-private-field.vue';

export default {
  inject: ['$validator'],
  components: { ViewGroupField, ViewTagsField, ViewDuplicatePrivateField, PeriodicRefreshField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    groups: {
      type: Array,
      default: () => [],
    },
    duplicatePrivate: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
