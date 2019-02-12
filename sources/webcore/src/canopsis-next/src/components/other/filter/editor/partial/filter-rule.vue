<template lang="pug">
  v-card.my-2.pa-0
    v-layout(justify-end)
      v-btn(
      @click="$emit('deleteRule')",
      color="red",
      small,
      flat,
      dark,
      fab
      )
        v-icon close
    v-layout.px-2(row, wrap, justify-space-around)
      v-flex.pa-1(xs12, md4)
        v-combobox.my-2(
        :items="possibleFields",
        :value="rule.field",
        @input="updateField('field', $event)",
        solo-inverted,
        hide-details,
        dense,
        flat
        )
      v-flex.pa-1(xs12, md4)
        v-combobox.my-2(
        :value="rule.operator",
        :items="operators",
        @input="updateField('operator', $event)",
        solo-inverted,
        hide-details,
        dense,
        flat
        )
      v-flex.pa-1(xs12, md4)
        v-layout.my-2(align-center, v-show="isShownInputField")
          v-flex(xs3)
            v-select(
            :items="types",
            :value="inputType",
            solo-inverted,
            hide-details,
            dense,
            flat,
            @input="updateInputTypeField"
            )
              template(slot="selection", slot-scope="{ parent, item, index }")
                v-icon.type-icon(small) {{ getIcon(item.value) }}
              template(slot="item", slot-scope="{ item }")
                v-list-tile-avatar.small-avatar
                  v-icon.type-icon(small) {{ getIcon(item.value) }}
                v-list-tile-content
                  v-list-tile-title {{ item.text }}
          v-flex(xs9)
            v-text-field.input-field(
            v-show="isInputType",
            :type="rule.inputType === $constants.FILTER_INPUT_TYPES.number ? 'number' : 'text'",
            :value="rule.input",
            @input="updateInputField",
            solo-inverted,
            hide-details,
            single-line,
            flat
            )
            v-switch.ma-0.ml-3.pa-0(
            v-show="$constants.FILTER_INPUT_TYPES.boolean === inputType",
            :inputValue="rule.input",
            @change="updateField('input', $event)",
            solo-inverted,
            hide-details,
            single-line,
            flat
            )
</template>

<script>
import { FILTER_OPERATORS, FILTER_INPUT_TYPES } from '@/constants';

import { getInputType } from '@/helpers/filter/editor/parse-group-to-filter';

import formMixin from '@/mixins/form';

/**
 * Component representing a rule in MongoDB filter
 *
 * @prop {Object} rule - Object of the rule
 * @prop {Array} possibleFields - List of all possible fields to filter on
 * @prop {Array} [operators=Object.values(FILTER_OPERATORS)] - List of all possible operators. Ex : 'equal', ...
 *
 * @event field#update
 * @event operator#update
 * @event input#update
 * @event deleteRule#click
 */
export default {
  mixins: [formMixin],
  model: {
    prop: 'rule',
    event: 'update:rule',
  },
  props: {
    rule: {
      type: Object,
      required: true,
    },
    possibleFields: {
      type: Array,
      required: true,
    },
    operators: {
      type: Array,
      default() {
        return Object.values(FILTER_OPERATORS);
      },
    },
  },
  data() {
    return {
      types: [
        { text: 'String', value: FILTER_INPUT_TYPES.string },
        { text: 'Number', value: FILTER_INPUT_TYPES.number },
        { text: 'Boolean', value: FILTER_INPUT_TYPES.boolean },
      ],
    };
  },
  computed: {
    getIcon() {
      const TYPES_ICONS_MAP = {
        [FILTER_INPUT_TYPES.string]: 'title',
        [FILTER_INPUT_TYPES.number]: 'exposure_plus_1',
        [FILTER_INPUT_TYPES.boolean]: 'toggle_on',
      };

      return type => TYPES_ICONS_MAP[type];
    },
    inputType() {
      return getInputType(this.rule.input);
    },
    isInputType() {
      return [FILTER_INPUT_TYPES.number, FILTER_INPUT_TYPES.string].includes(this.inputType);
    },
    isShownInputField() {
      return ![
        FILTER_OPERATORS.isEmpty,
        FILTER_OPERATORS.isNotEmpty,
        FILTER_OPERATORS.isNull,
        FILTER_OPERATORS.isNotNull,
      ].includes(this.rule.operator);
    },
  },
  methods: {
    updateInputField(value) {
      const isInputTypeNumber = this.inputType === FILTER_INPUT_TYPES.number;

      this.updateField('input', isInputTypeNumber ? Number(value) : value);
    },
    updateInputTypeField(value) {
      switch (value) {
        case FILTER_INPUT_TYPES.number:
          this.updateField('input', Number(this.rule.input));
          break;
        case FILTER_INPUT_TYPES.boolean:
          this.updateField('input', Boolean(this.rule.input));
          break;
        case FILTER_INPUT_TYPES.string:
          this.updateField('input', String(this.rule.input));
          break;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .input-field {
    border-left: solid 1px #cccccc;
  }

  .type-icon {
    color: inherit;
    opacity: .6;
  }

  .small-avatar {
    min-width: 30px;

    & /deep/ .v-avatar {
      width: 20px!important;
      height: 20px!important;
    }
  }
</style>
