<template>
  <v-layout>
    <v-flex
      class="pr-2"
      xs3
    >
      <v-select
        v-field="condition.type"
        :items="conditionTypes"
        :label="$t('common.type')"
        :disabled="disabled"
      />
    </v-flex>
    <v-flex
      class="px-2"
      xs4
    >
      <v-text-field
        v-field="condition.attribute"
        v-validate="'required'"
        :label="$t('common.attribute')"
        :name="conditionFieldName"
        :error-messages="errors.collect(conditionFieldName)"
        :disabled="disabled"
      />
    </v-flex>
    <v-flex
      class="pl-2"
      xs5
    >
      <v-layout align-center>
        <c-payload-text-field
          v-field="condition.value"
          :label="$t('common.value')"
          :disabled="disabled"
          :variables="variables"
          :name="valueFieldName"
          clearable
        />
        <v-btn
          v-if="!disabled"
          :disabled="disabledRemove"
          icon
          small
          @click="removeCondition"
        >
          <v-icon
            color="error"
            small
          >
            delete
          </v-icon>
        </v-btn>
      </v-layout>
    </v-flex>
  </v-layout>
</template>

<script>
import { EXTERNAL_DATA_CONDITION_TYPES } from '@/constants';

import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'condition',
    event: 'input',
  },
  props: {
    condition: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    disabledRemove: {
      type: Boolean,
      default: false,
    },
    variables: {
      type: Array,
      default: () => ([]),
    },
  },
  computed: {
    conditionTypes() {
      return Object.values(EXTERNAL_DATA_CONDITION_TYPES)
        .map(type => ({ text: this.$t(`externalData.conditionTypes.${type}`), value: type }));
    },

    conditionFieldName() {
      return `${this.name}.condition`;
    },

    valueFieldName() {
      return `${this.name}.value`;
    },
  },
  methods: {
    removeCondition() {
      this.$emit('remove', this.condition);
    },
  },
};
</script>
