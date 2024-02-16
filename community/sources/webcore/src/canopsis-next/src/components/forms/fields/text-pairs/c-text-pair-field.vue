<template>
  <v-layout align-center>
    <v-layout>
      <v-flex
        class="mr-3"
        xs6
      >
        <v-text-field
          v-field="item[itemText]"
          v-validate="textValidationRules"
          :label="textLabel"
          :disabled="disabled"
          :name="textFieldName"
          :error-messages="errors.collect(textFieldName)"
        />
      </v-flex>
      <v-flex xs6>
        <v-text-field
          v-field="item[itemValue]"
          v-validate="valueValidationRules"
          :label="valueLabel"
          :disabled="disabled"
          :name="valueFieldName"
          :error-messages="errors.collect(valueFieldName)"
        >
          <template #append="">
            <slot name="append-value" />
          </template>
        </v-text-field>
      </v-flex>
    </v-layout>
    <c-action-btn
      v-if="!disabled"
      type="delete"
      @click="$emit('remove')"
    />
  </v-layout>
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'item',
    event: 'input',
  },
  props: {
    item: {
      type: Object,
      default: () => ({}),
    },
    textLabel: {
      type: String,
      default: '',
    },
    valueLabel: {
      type: String,
      default: '',
    },
    itemText: {
      type: String,
      default: 'text',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    name: {
      type: String,
      default: 'item',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    textRequired: {
      type: Boolean,
      default: false,
    },
    valueRequired: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    textValidationRules() {
      return {
        required: this.textRequired,
      };
    },

    valueValidationRules() {
      return {
        required: this.valueRequired,
      };
    },

    textFieldName() {
      return `${this.name}.${this.itemText}`;
    },

    valueFieldName() {
      return `${this.name}.${this.itemValue}`;
    },
  },
};
</script>
