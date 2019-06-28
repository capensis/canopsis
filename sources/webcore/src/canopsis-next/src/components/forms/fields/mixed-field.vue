<template lang="pug">
  .mixed-field(:class="{ 'mixed-field__solo-inverted': soloInverted }")
    v-select.mixed-field__type-selector(
    :items="types",
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
    v-text-field.mixed-field__text(
    v-if="isInputTypeText",
    :type="inputType === $constants.FILTER_INPUT_TYPES.number ? 'number' : 'text'",
    :value="value",
    :name="name",
    :disabled="disabled",
    :solo-inverted="soloInverted",
    :hide-details="hideDetails",
    :flat="flat",
    :error-messages="errorMessages",
    single-line,
    dense,
    @input="updateTextFieldValue",
    )
    v-switch.ma-0.ml-3.mixed-field__switch(
    v-else-if="inputType === $constants.FILTER_INPUT_TYPES.boolean",
    :inputValue="value",
    :label="switchLabel",
    :disabled="disabled",
    hide-details,
    @change="updateModel"
    )
    v-text-field.mixed-field__text(
    v-else,
    :error-messages="errorMessages",
    value="null",
    disabled
    )
</template>

<script>
import { isBoolean, isNumber, isNan, isNull } from 'lodash';

import { FILTER_INPUT_TYPES } from '@/constants';

import formBaseMixin from '@/mixins/form/base';

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
      type: [String, Number, Boolean],
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
  },
  data() {
    return {
      types: [
        { text: 'String', value: FILTER_INPUT_TYPES.string },
        { text: 'Number', value: FILTER_INPUT_TYPES.number },
        { text: 'Boolean', value: FILTER_INPUT_TYPES.boolean },
        { text: 'Null', value: FILTER_INPUT_TYPES.null },
      ],
    };
  },
  computed: {
    switchLabel() {
      return String(this.value);
    },

    inputType() {
      if (isBoolean(this.value)) {
        return FILTER_INPUT_TYPES.boolean;
      } else if (isNumber(this.value)) {
        return FILTER_INPUT_TYPES.number;
      } else if (isNull(this.value)) {
        return FILTER_INPUT_TYPES.null;
      }

      return FILTER_INPUT_TYPES.string;
    },

    isInputTypeText() {
      return [FILTER_INPUT_TYPES.number, FILTER_INPUT_TYPES.string].includes(this.inputType);
    },

    getInputTypeIcon() {
      const TYPES_ICONS_MAP = {
        [FILTER_INPUT_TYPES.string]: 'title',
        [FILTER_INPUT_TYPES.number]: 'exposure_plus_1',
        [FILTER_INPUT_TYPES.boolean]: '$vuetify.icons.toggle_on',
        [FILTER_INPUT_TYPES.null]: 'space_bar',
      };

      return type => TYPES_ICONS_MAP[type];
    },
  },
  methods: {
    updateType(value) {
      switch (value) {
        case FILTER_INPUT_TYPES.number:
          this.updateModel(Number(this.value));
          break;
        case FILTER_INPUT_TYPES.boolean:
          this.updateModel(Boolean(this.value));
          break;
        case FILTER_INPUT_TYPES.string:
          this.updateModel((isNan(this.value) || isNull(this.value)) ? '' : String(this.value));
          break;
        case FILTER_INPUT_TYPES.null:
          this.updateModel(null);
          break;
      }
    },

    updateTextFieldValue(value) {
      const isInputTypeNumber = this.inputType === FILTER_INPUT_TYPES.number;

      this.updateModel(isInputTypeNumber ? Number(value) : value);
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
          width: 20px!important;
          height: 20px!important;
        }
      }
    }
  }
</style>
