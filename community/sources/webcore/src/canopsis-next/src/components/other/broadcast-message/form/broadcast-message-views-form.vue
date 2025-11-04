<template>
  <v-layout class="gap-2 py-2" column>
    <v-treeview
      v-field="views"
      :items="treeItems"
      selected-color="primary"
      item-key="value"
      item-children="children"
      selectable
      open-on-click
      transition
      @input="asyncValidateRequiredRule"
    >
      <template #label="{ item }">
        <span>{{ item.name }}</span>
      </template>
    </v-treeview>
    <c-alert :value="errors.has('views')" type="error">
      {{ $t('broadcastMessage.errors.viewsRequired') }}
    </c-alert>
  </v-layout>
</template>

<script>
import { useValidationAttachRequiredForField } from '@/hooks/validator/validation-attach-required';

export default {
  inject: ['$validator'],
  model: {
    prop: 'views',
    event: 'input',
  },
  props: {
    views: {
      type: Array,
      default: () => [],
    },
    treeItems: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const { asyncValidateRequiredRule } = useValidationAttachRequiredForField('views', () => props.views.length > 0);

    return {
      asyncValidateRequiredRule,
    };
  },
};
</script>
