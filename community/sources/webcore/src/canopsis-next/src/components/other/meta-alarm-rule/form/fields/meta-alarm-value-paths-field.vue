<template>
  <v-layout
    class="meta-alarm-value-paths-field"
    wrap
  >
    <v-flex xs12>
      <v-layout
        v-for="(item, index) in items"
        :key="item[itemKey]"
        class="text-field"
        justify-space-between
        align-center
      >
        <v-flex xs12>
          <v-text-field
            v-validate="validationRules"
            :value="item[itemValue]"
            :label="label"
            :name="getFieldName(item[itemKey])"
            :error-messages="errors.collect(getFieldName(item[itemKey]))"
            @input="updateFieldInArrayItem(index, itemValue, $event)"
          />
        </v-flex>
        <c-action-btn type="delete" @click="removeItemFromArray(index)" />
        <c-help-icon :text="$t('metaAlarmRule.valuePathHelpText')" icon="help" top />
      </v-layout>
    </v-flex>
    <v-flex xs12>
      <v-layout class="gap-2" align-center>
        <v-btn
          :color="error ? 'error' : 'primary'"
          @click="addNewItem"
        >
          {{ $t('common.add') }}
        </v-btn>
        <span
          v-show="error"
          class="error--text"
        >
          {{ $t('metaAlarmRule.errors.noValuePaths') }}
        </span>
      </v-layout>
    </v-flex>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { defaultPrimitiveArrayItemCreator } from '@/helpers/entities/shared/form';

import CHelpIcon from '@/components/common/icons/c-help-icon.vue';

import { useArrayModelField } from '@/hooks/form/useArrayModelField';
import { useInjectValidator } from '@/hooks/form/useValidationChildren';

export default {
  inject: ['$validator'],
  components: { CHelpIcon },
  model: {
    prop: 'items',
    event: 'input',
  },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    itemKey: {
      type: String,
      default: 'key',
    },
    name: {
      type: String,
      default: 'item',
    },
    required: {
      type: Boolean,
      default: false,
    },
    error: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const validator = useInjectValidator();
    const { addItemIntoArray, updateFieldInArrayItem, removeItemFromArray } = useArrayModelField(props, emit);

    const validationRules = computed(() => ({
      required: props.required,
    }));

    const getNamePrefix = key => `${props.name}[${key}]`;
    const getFieldName = key => `${getNamePrefix(key)}.${props.itemValue}`;

    const addNewItem = () => {
      addItemIntoArray(defaultPrimitiveArrayItemCreator());
    };

    return {
      errors: validator.errors,
      validationRules,
      addNewItem,
      getFieldName,
      updateFieldInArrayItem,
      removeItemFromArray,
    };
  },
};
</script>

<style lang="scss" scoped>
.meta-alarm-value-paths-field {
}
</style>
