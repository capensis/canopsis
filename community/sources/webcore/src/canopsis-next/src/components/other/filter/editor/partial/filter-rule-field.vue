<template lang="pug">
  v-combobox(
    ref="combobox",
    :value="preparedField",
    :items="preparedItems",
    :value-preparer="valuePreparer",
    :return-object="false",
    item-text="name",
    item-value="value",
    solo-inverted,
    hide-details,
    dense,
    flat,
    @input="change",
    @update:searchInput="change"
  )
    template(slot="item", slot-scope="props")
      v-list-tile-content {{ props.item | itemText }}
    template(slot="prepend-inner", v-if="selectedField.additionalFieldProps")
      v-chip(
        color="grey lighten-1",
        close,
        @input="remove"
      ) {{ selectedField.name }}
</template>

<script>
import { get, isObject, isString, isUndefined } from 'lodash';

import { formBaseMixin } from '@/mixins/form';

export default {
  filters: {
    itemText(item) {
      if (isObject(item)) {
        return `${item.name} (${item.value})`;
      }

      return item;
    },
  },
  mixins: [formBaseMixin],
  model: {
    prop: 'field',
    event: 'input',
  },
  props: {
    field: {
      type: String,
      default: '',
    },
    items: {
      type: Array,
      default: () => [],
    },
    selectedField: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      searchInput: '',
    };
  },
  computed: {
    isAdditionalField() {
      return !!this.selectedField.additionalFieldProps;
    },

    isAdditionalFieldLastPath() {
      const [,, additionalField] = this.field.split('.');

      return !isUndefined(additionalField);
    },

    preparedField() {
      if (this.isAdditionalField) {
        return this.field.replace(new RegExp(`^${this.selectedField.value}.?`), '');
      }

      return this.field;
    },

    preparedItems() {
      if (this.isAdditionalField && !this.isAdditionalFieldLastPath) {
        return [];
      }

      return get(this.selectedField, 'additionalFieldProps.items', this.items);
    },
  },
  watch: {
    selectedField({ value, additionalFieldProps } = {}, { value: oldValue } = {}) {
      if (value !== oldValue && additionalFieldProps) {
        /**
         * We are clearing internalValue in combobox when new selectedField was selected
         */
        this.$nextTick(() => {
          this.$refs.combobox.internalValue = '';
          this.$refs.combobox.focus();
        });
      }
    },
  },
  methods: {
    valuePreparer(value) {
      if (isString(value) && this.isAdditionalField) {
        const [, parent] = this.field.split('.');
        return `${parent}.${value}`;
      }

      return value;
    },

    remove() {
      this.updateModel('');
    },

    change(field = '') {
      let preparedField = field || '';

      if (this.isAdditionalField) {
        preparedField = `${this.selectedField.value}.${preparedField}`;
      }

      this.updateModel(preparedField);
    },
  },
};
</script>
