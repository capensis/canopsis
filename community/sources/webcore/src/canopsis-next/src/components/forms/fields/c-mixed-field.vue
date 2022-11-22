<template lang="pug">
  div.mixed-field(:class="{ 'mixed-field__solo-inverted': soloInverted }")
    v-select.mixed-field__type-selector(
      :items="preparedTypes",
      :value="inputType",
      :label="label",
      :disabled="disabled",
      :solo-inverted="soloInverted",
      :flat="flat",
      :error-messages="errorMessages",
      hide-details,
      dense,
      @input="updateType"
    )
      template(slot="selection", slot-scope="{ parent, item, index }")
        v-icon.mixed-field__type-selector-icon(small) {{ getInputTypeIcon(item.value) }}
      template(slot="item", slot-scope="{ item }")
        v-list-tile-avatar.mixed-field__type-selector-avatar
          v-icon.mixed-field__type-selector-icon(small) {{ getInputTypeIcon(item.value) }}
        v-list-tile-content
          v-list-tile-title {{ item.text }}
    component(
      :is="inputComponent.is",
      v-bind="inputComponent.bind",
      v-on="inputComponent.on"
    )
</template>

<script>
import { isBoolean, isNumber, isNan, isArray, isUndefined, isEmpty, isNull, pick } from 'lodash';

import { FILTER_INPUT_TYPES } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

export default {
  $_veeValidate: {
    name() {
      return this.name;
    },

    value() {
      return this.value;
    },
  },
  inject: ['$validator'],
  mixins: [formBaseMixin],
  props: {
    value: {
      type: [String, Number, Boolean, Array],
      default: '',
    },
    name: {
      type: String,
      default: null,
    },
    label: {
      type: String,
      default: null,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    soloInverted: {
      type: Boolean,
      default: false,
    },
    flat: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
    items: {
      type: Array,
      default: () => [],
    },
    itemText: {
      type: String,
      default: 'text',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    types: {
      type: Array,
      default: () => [
        { value: FILTER_INPUT_TYPES.string },
        { value: FILTER_INPUT_TYPES.number },
        { value: FILTER_INPUT_TYPES.boolean },
        { value: FILTER_INPUT_TYPES.null },
        { value: FILTER_INPUT_TYPES.array },
      ],
    },
  },
  computed: {
    preparedTypes() {
      return this.types.map(
        type => (type.text ? type : ({ ...type, text: this.$t(`mixedField.types.${type.value}`) })),
      );
    },

    switchLabel() {
      return String(this.value);
    },

    inputType() {
      if (isBoolean(this.value)) {
        return FILTER_INPUT_TYPES.boolean;
      }

      if (isNumber(this.value)) {
        return FILTER_INPUT_TYPES.number;
      }

      if (isNull(this.value)) {
        return FILTER_INPUT_TYPES.null;
      }

      if (isArray(this.value)) {
        return FILTER_INPUT_TYPES.array;
      }

      return FILTER_INPUT_TYPES.string;
    },

    isInputTypeText() {
      return [FILTER_INPUT_TYPES.number, FILTER_INPUT_TYPES.string].includes(this.inputType);
    },

    inputComponent() {
      if (this.isInputTypeText) {
        const additionalProps = this.items.length
          ? { ...pick(this, ['items', 'itemText', 'itemValue']), returnObject: false, forceSearching: true }
          : {};

        return {
          is: this.items.length ? 'v-combobox' : 'v-text-field',
          bind: {
            ...pick(this, [
              'value',
              'name',
              'disabled',
              'soloInverted',
              'hideDetails',
              'flat',
              'errorMessages',
            ]),
            ...additionalProps,

            class: 'mixed-field__text',
            type: this.inputType === FILTER_INPUT_TYPES.number ? 'number' : 'text',
            singleLine: true,
            dense: true,
          },
          on: {
            input: this.updateTextFieldValue,
            'update:searchInput': this.updateTextFieldValue,
          },
        };
      }

      if (this.inputType === FILTER_INPUT_TYPES.boolean) {
        return {
          is: 'v-switch',

          bind: {
            class: 'ma-0 ml-3 mixed-field__switch',
            inputValue: this.value,
            label: this.switchLabel,
            disabled: this.disabled,
            color: 'primary',
            hideDetails: true,
          },
          on: {
            change: this.updateModel,
          },
        };
      }

      if (this.inputType === FILTER_INPUT_TYPES.array) {
        return {
          is: 'c-array-mixed-field',

          bind: {
            class: 'ml-3',
            values: this.value,
            disabled: this.disabled,
            types: this.types.filter(({ value }) => value !== FILTER_INPUT_TYPES.array),
          },
          on: {
            change: this.updateModel,
          },
        };
      }

      return {
        is: 'v-text-field',

        bind: {
          class: 'mixed-field__text',
          errorMessages: this.errorMessages,
          value: 'null',
          disabled: true,
        },
      };
    },
  },
  watch: {
    types: {
      immediate: true,
      handler(types) {
        if (!types.some(({ value }) => value === this.inputType)) {
          const [type = {}] = types;

          this.updateType(type.value, type.defaultValue);
        }
      },
    },
  },
  methods: {
    getInputTypeIcon(type) {
      return {
        [FILTER_INPUT_TYPES.string]: 'title',
        [FILTER_INPUT_TYPES.number]: 'exposure_plus_1',
        [FILTER_INPUT_TYPES.boolean]: 'toggle_on',
        [FILTER_INPUT_TYPES.null]: 'space_bar',
        [FILTER_INPUT_TYPES.array]: 'view_array',
      }[type];
    },

    updateType(type, defaultValueForType) {
      const preparedValue = isEmpty(this.value) && !isUndefined(defaultValueForType)
        ? defaultValueForType
        : this.value;

      switch (type) {
        case FILTER_INPUT_TYPES.number:
          this.updateModel(Number(preparedValue));
          break;
        case FILTER_INPUT_TYPES.boolean:
          this.updateModel(Boolean(preparedValue));
          break;
        case FILTER_INPUT_TYPES.string:
          this.updateModel((isNan(preparedValue) || isNull(preparedValue)) ? '' : String(preparedValue));
          break;
        case FILTER_INPUT_TYPES.null:
          this.updateModel(null);
          break;
        case FILTER_INPUT_TYPES.array:
          this.updateModel(preparedValue ? [preparedValue] : []);
          break;
        default:
          this.updateModel(undefined);
      }
    },

    updateTextFieldValue(value) {
      let preparedValue = value;

      if (isNull(value) && this.inputType !== FILTER_INPUT_TYPES.null) {
        preparedValue = '';
      }

      if (this.inputType === FILTER_INPUT_TYPES.number) {
        preparedValue = Number(value);
      }

      this.updateModel(preparedValue);
    },
  },
};
</script>

<style lang="scss" scoped>
  .mixed-field {
    position: relative;
    padding-left: 45px;

    &__solo-inverted {
      padding-left: 60px;

      &.mixed-field__switch {
      }
    }

    &__text {
      margin-top: 0;

      & /deep/ input {
        padding-left: 5px;
      }

      & /deep/ .v-text-field__details {
        margin-left: -45px;

        .mixed-field__solo-inverted & {
          margin-left: -60px;
        }
      }
    }

    &__switch {
      padding: 18px 0;

      & /deep/ .v-label {
        text-transform: capitalize;
      }

      .mixed-field__solo-inverted & {
        padding: 12px 0;
      }
    }

    &__type-selector {
      width: 45px;
      position: absolute;
      margin: 0;
      left: 0;
      top: 0;

      &.v-text-field--solo-inverted {
        width: 59px;
      }

      & /deep/ .v-input__slot {
        padding-right: 5px;
      }

      &-icon {
        color: inherit;
        opacity: .6;
      }

      &-avatar {
        min-width: 30px;

        & /deep/ .v-avatar {
          width: 20px !important;
          height: 20px !important;
        }
      }
    }
  }
</style>
