<template>
  <v-layout column>
    <c-draggable-list-field
      v-field="columns"
      :class="{ empty: isColumnsEmpty }"
      :handle="`.${dragItemHandleClass}`"
    >
      <column-field
        v-for="(column, index) in columns"
        v-field="columns[index]"
        :key="column.key"
        :name="column.key"
        :type="type"
        :drag-handle-class="dragItemHandleClass"
        :with-html="withHtml"
        :with-template="withTemplate"
        :with-color-indicator="withColorIndicator"
        :with-instructions="withInstructions"
        :optional-infos-attributes="optionalInfosAttributes"
        :with-simple-template="withSimpleTemplate"
        :without-infos-attributes="withoutInfosAttributes"
        :variables="variables"
        :excluded-columns="excludedColumns"
        class="mb-3"
        @remove="remove(index)"
      />
    </c-draggable-list-field>
    <v-layout justify-end>
      <v-tooltip left>
        <template #activator="{ on }">
          <v-btn
            class="mr-2 mx-0"
            color="primary"
            fab
            small
            v-on="on"
            @click.prevent="add"
          >
            <v-icon>add</v-icon>
          </v-btn>
        </template>
        <span>{{ $t('common.add') }}</span>
      </v-tooltip>
    </v-layout>
  </v-layout>
</template>

<script>
import { ENTITIES_TYPES } from '@/constants';

import { widgetColumnToForm } from '@/helpers/entities/widget/column/form';

import { formArrayMixin, formValidationHeaderMixin } from '@/mixins/form';

import ColumnField from './partials/column-field.vue';

export default {
  inject: ['$validator'],
  components: { ColumnField },
  mixins: [
    formArrayMixin,
    formValidationHeaderMixin,
  ],
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    type: {
      type: String,
      default: ENTITIES_TYPES.alarm,
    },
    columns: {
      type: [Array, Object],
      default: () => [],
    },
    withTemplate: {
      type: Boolean,
      default: false,
    },
    withHtml: {
      type: Boolean,
      default: false,
    },
    withColorIndicator: {
      type: Boolean,
      default: false,
    },
    withInstructions: {
      type: Boolean,
      default: false,
    },
    optionalInfosAttributes: {
      type: Boolean,
      default: false,
    },
    withSimpleTemplate: {
      type: Boolean,
      default: false,
    },
    variables: {
      type: Array,
      required: false,
    },
    withoutInfosAttributes: {
      type: Boolean,
      default: false,
    },
    excludedColumns: {
      type: Array,
      required: false,
    },
  },
  computed: {
    dragItemHandleClass() {
      return 'column-drag-handle';
    },

    isColumnsEmpty() {
      return !this.columns.length;
    },
  },
  methods: {
    add() {
      this.addItemIntoArray(widgetColumnToForm());
    },

    remove(index) {
      this.removeItemFromArray(index);
    },
  },
};
</script>
