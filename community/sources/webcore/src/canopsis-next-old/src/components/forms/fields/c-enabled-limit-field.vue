<template>
  <v-layout wrap>
    <v-flex xs5>
      <v-checkbox
        v-validate
        v-field="value.enabled"
        :label="label"
        :error-messages="errors.collect(enabledFieldName)"
        :name="enabledFieldName"
        color="primary"
      >
        <template #append="">
          <c-help-icon
            v-if="helpText"
            :text="helpText"
            max-width="300"
            color="info"
            top
          />
        </template>
      </v-checkbox>
    </v-flex>
    <v-flex xs2>
      <c-number-field
        v-field="value.limit"
        v-bind="validationRules"
        :label="fieldLabel"
        :name="limitFieldName"
        :disabled="!value.enabled"
      />
    </v-flex>
    <v-flex xs9>
      <v-messages
        :value="errors.collect(name)"
        color="error"
      />
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    fieldLabel: {
      type: String,
      default: '',
    },
    helpText: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'limit',
    },
    min: {
      type: Number,
      required: false,
    },
  },
  computed: {
    enabledFieldName() {
      return `${this.name}.enabled`;
    },

    limitFieldName() {
      return `${this.name}.limit`;
    },

    validationRules() {
      const rules = {};

      if (this.value.enabled) {
        rules.min = this.min;
        rules.reqired = true;
      }

      return rules;
    },
  },
};
</script>
