<template>
  <v-layout
    :class="{ 'text-fields__disabled': disabled }"
    class="text-fields"
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
            :disabled="disabled"
            :name="getFieldName(item[itemKey])"
            :error-messages="errors.collect(getFieldName(item[itemKey]))"
            @input="updateFieldInArrayItem(index, itemValue, $event)"
          />
        </v-flex>
        <div class="text-fields__delete-button">
          <v-btn
            v-if="!disabled"
            icon
            @click="removeItemFromArray(index)"
          >
            <v-icon color="error">
              delete
            </v-icon>
          </v-btn>
        </div>
      </v-layout>
    </v-flex>
    <v-flex
      v-if="!disabled"
      xs12
    >
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

import { formArrayMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin],
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
    disabled: {
      type: Boolean,
      default: false,
    },
    error: {
      type: Boolean,
      default: false,
    },
  },
  computed(props) {
    const validationRules = computed(() => ({
      required: props.required,
    }));

    return {
      validationRules,
    };
  },
  methods: {
    getNamePrefix(key) {
      return `${this.name}[${key}]`;
    },

    getTextFieldName(key) {
      return `${this.getNamePrefix(key)}.${this.itemText}`;
    },

    getFieldName(key) {
      return `${this.getNamePrefix(key)}.${this.itemValue}`;
    },

    addNewItem() {
      this.addItemIntoArray(defaultPrimitiveArrayItemCreator());
    },
  },
};
</script>

<style lang="scss" scoped>
.text-fields {
  &:not(.text-fields__disabled) .text-pair {
    position: relative;
    padding-right: 50px;

    &__delete-button {
      position: absolute;
      right: 0;
      top: 50%;
      transform: translateY(-50%);
    }
  }
}
</style>
