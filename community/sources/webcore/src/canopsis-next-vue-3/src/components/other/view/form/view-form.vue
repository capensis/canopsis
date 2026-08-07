<template>
  <c-layout column>
    <view-duplicate-private-field
      v-if="duplicatePrivate"
      :value="form.is_private"
      @update:value="updateField('is_private', $event)"
    />

    <c-name-field
      :model-value="form.title"
      :label="$t('common.title')"
      name="title"
      autofocus
      required
      @update:model-value="updateField('title', $event)"
    />

    <c-text-field
      :model-value="form.description"
      :label="$t('common.description')"
      name="description"
      @update:model-value="updateField('description', $event)"
    />

    <c-enabled-field
      :model-value="form.enabled"
      @update:model-value="updateField('enabled', $event)"
    />

    <periodic-refresh-field
      :periodic-refresh="form.periodic_refresh"
      :label="$t('view.periodicRefresh')"
      @update:periodic-refresh="updateField('periodic_refresh', $event)"
    />

    <view-tags-field
      :value="form.tags"
      @update:value="updateField('tags', $event)"
    />

    <view-group-field
      :value="form.group"
      :groups="groups"
      @update:value="updateField('group', $event)"
    />
  </c-layout>
</template>

<script>
import { useModelField } from '@/hooks/form';

import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';

import ViewTagsField from './fields/view-tags-field.vue';
import ViewGroupField from './fields/view-group-field.vue';
import ViewDuplicatePrivateField from './fields/view-duplicate-private-field.vue';

export default {
  components: {
    ViewGroupField,
    ViewTagsField,
    ViewDuplicatePrivateField,
    PeriodicRefreshField,
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

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit, { modelProp: 'form' });

    return {
      updateField,
    };
  },
};
</script>
