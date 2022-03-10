<template lang="pug">
  v-layout(row)
    v-flex(xs4)
      v-layout(row)
        v-flex(:xs4="hasAdditionalAttributeField")
          c-pattern-attribute-field(v-field="value.attribute", :items="attributes")
        v-flex.pl-2(v-if="hasAdditionalAttributeField", xs8)
          component(
            v-field="value",
            v-bind="additionalAttributeField.props",
            v-on="additionalAttributeField.on",
            :is="additionalAttributeField.is"
          )
    v-flex.pl-2(xs4)
      component(
        v-if="operatorField.is",
        v-field="value.operator",
        v-bind="operatorField.props",
        v-on="operatorField.on",
        :is="operatorField.is"
      )
      c-pattern-operator-field(
        v-else,
        v-field="value.operator",
        v-bind="operatorField.props",
        v-on="operatorField.on",
        :operators="operators"
      )
    v-flex.pl-2(xs4)
      component(
        v-if="valueField.is",
        v-field="value.value",
        v-bind="valueField.props",
        v-on="valueField.on",
        :is="valueField.is"
      )
      v-text-field(
        v-else,
        v-field="value.value",
        v-bind="valueField.props",
        v-on="valueField.on",
        :label="$t('common.value')"
      )
</template>

<script>
import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    attributes: {
      type: Array,
      default: () => [],
    },
    operators: {
      type: Array,
      default: () => [],
    },
    operatorField: {
      type: Object,
      default: () => ({}),
    },
    valueField: {
      type: Object,
      default: () => ({}),
    },
    additionalAttributeField: {
      type: Object,
      required: false,
    },
    name: {
      type: String,
      default: 'rule',
    },
  },
  computed: {
    hasAdditionalAttributeField() {
      return !!this.additionalAttributeField;
    },
  },
};
</script>
