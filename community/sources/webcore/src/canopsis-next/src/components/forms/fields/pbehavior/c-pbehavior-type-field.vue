<template>
  <v-select
    v-field="value"
    v-validate="rules"
    :label="label || $t('common.type')"
    :loading="fieldPbehaviorTypesPending"
    :items="preparedItems"
    :error-messages="errors.collect(name)"
    :name="name"
    :disabled="disabled"
    :multiple="multiple"
    :chips="chips"
    :deletable-chips="chips"
    :small-chips="chips"
    :item-disabled="isItemDisabled"
    :return-object="returnObject"
    item-text="name"
    item-value="_id"
  />
</template>

<script>
import { isArray, isObject, isEmpty } from 'lodash';

import { mapIds } from '@/helpers/array';

import { entitiesFieldPbehaviorFieldTypeMixin } from '@/mixins/entities/pbehavior/types-field';

export default {
  inject: ['$validator'],
  mixins: [entitiesFieldPbehaviorFieldTypeMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Object, String, Array],
      default: '',
    },
    name: {
      type: String,
      default: 'type',
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    chips: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    max: {
      type: Number,
      required: false,
    },
  },
  data() {
    return {
      originalValue: this.value,
    };
  },
  computed: {
    selectedTypesIds() {
      return this.getSelectedTypesIds();
    },

    originalSelectedTypesIds() {
      return this.getSelectedTypesIds(this.originalValue);
    },

    preparedItems() {
      return this.fieldPbehaviorTypes.filter(type => (
        !type.hidden || this.originalSelectedTypesIds.includes(type._id)
      ));
    },

    rules() {
      return {
        required: this.required,
      };
    },
  },
  methods: {
    isItemDisabled(item) {
      if (this.max) {
        return this.value.length === this.max && !this.selectedTypesIds.includes(item._id);
      }

      return false;
    },

    getSelectedTypesIds(value = this.value) {
      if (isArray(value)) {
        return this.returnObject
          ? mapIds(value)
          : value;
      }

      return isEmpty(value)
        ? []
        : [
          isObject(value)
            ? value._id
            : value,
        ];
    },
  },
};
</script>
