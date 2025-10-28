<template>
  <v-combobox
    v-validate="rules"
    ref="field"
    :value="value"
    :search-input="value"
    :label="label ?? $t('common.payload')"
    :items="variables"
    :disabled="disabled"
    :return-object="false"
    :menu-props="menuProps"
    :error-messages="errorMessages"
    :clearable="clearable"
    :name="name"
    no-filter
    @blur="handleBlur"
    @update:search-input="onSearchInputChange"
  >
    <template #append="">
      <slot name="append" />
    </template>
    <template #append-outer="">
      <slot name="append-outer" />
    </template>
    <template #list="">
      <variables-list
        :value="variablesMenuValue"
        :items="variables"
        children-key="variables"
        show-value
        hide-empty-value
        clickable-parent
        @input="pasteVariable"
      />
    </template>
  </v-combobox>
</template>

<script>
import { payloadFieldMixin } from '@/mixins/payload/payload-field';

import VariablesList from '@/components/common/text-editor/variables-list.vue';

export default {
  inject: ['$validator'],
  components: { VariablesList },
  mixins: [payloadFieldMixin],
  props: {
    name: {
      type: String,
      default: 'payload',
    },
    label: {
      type: String,
      required: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

    hasValidationRules() {
      return Object.values(this.rules).some(rule => !!rule);
    },

    errorMessages() {
      return this.errors.collect(this.name).map((message = '') => (
        message.includes('|') ? message.split('|')[1] : message
      ));
    },

    menuProps() {
      return {
        value: !!this.variables.length && this.variablesShown,
        offsetY: true,
        closeOnContentClick: false,
        minWidth: 200,
      };
    },
  },
  methods: {
    onSearchInputChange(value) {
      this.debouncedOnSelectionChange();
      this.updateModel(value ?? '');

      if (this.errorMessages?.length) {
        this.$nextTick(() => {
          if (this.hasValidationRules) {
            this.$validator.validate(this.name);

            return;
          }

          this.errors.remove(this.name);
        });
      }
    },
  },
};
</script>
