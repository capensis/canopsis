<template lang="pug">
  v-layout(row)
    v-flex.pr-2(xs3)
      v-select(
        v-field="condition.type",
        :items="conditionTypes",
        :label="$t('common.type')"
      )
    v-flex.px-2(xs4)
      v-text-field(
        v-field="condition.attribute",
        v-validate="'required'",
        :label="$t('common.attribute')",
        :name="conditionFieldName",
        :error-messages="errors.collect(conditionFieldName)"
      )
    v-flex.pl-2(xs5)
      v-layout(row, align-center)
        v-combobox(
          ref="combobox",
          v-field="condition.value",
          :search-input.sync="searchInput",
          :label="$t('common.value')",
          :items="values",
          no-filter,
          clearable,
          @change="selectValue",
          @update:searchInput="debouncedOnSelectionChange"
        )
          template(#item="{ item, tile }")
            v-list-tile(:value="item.value === activeValue", @click="selectValue(item.value)")
              v-list-tile-content {{ item.text }}
              span.ml-4.grey--text {{ item.value }}
        v-btn(:disabled="disabledRemove", icon, small, @click="removeCondition(condition.key)")
          v-icon(color="red", small) delete
</template>

<script>
import { debounce } from 'lodash';

import {
  EVENT_FILTER_EXTERNAL_DATA_CONDITION_TYPES,
  EVENT_FILTER_EXTERNAL_DATA_CONDITION_VALUES,
} from '@/constants';

import { formMixin } from '@/mixins/form';
import { matchPayloadVariableBySelection } from '@/helpers/payload-json';

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
    disabledRemove: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      searchInput: this.condition.value,
      showItems: true,
      activeValue: undefined,
    };
  },
  computed: {
    values() {
      if (!this.showItems) {
        return [];
      }

      return Object.values(EVENT_FILTER_EXTERNAL_DATA_CONDITION_VALUES).map(({ value, text }) => ({
        value,
        text: this.$t(`eventFilter.externalDataValues.${text}`),
      }));
    },

    conditionTypes() {
      return Object.values(EVENT_FILTER_EXTERNAL_DATA_CONDITION_TYPES)
        .map(type => ({ text: this.$t(`eventFilter.externalDataConditionTypes.${type}`), value: type }));
    },

    conditionFieldName() {
      return `${this.name}.condition`;
    },
  },
  created() {
    this.debouncedOnSelectionChange = debounce(this.onSelectionChange, 50);
  },
  mounted() {
    document.addEventListener('selectionchange', this.debouncedOnSelectionChange);
  },
  beforeDestroy() {
    document.removeEventListener('selectionchange', this.debouncedOnSelectionChange);
  },
  methods: {
    selectValue(value) {
      const { selectionStart, selectionEnd } = this;

      this.activeValue = value ?? undefined;

      if (!this.searchInput) {
        this.updateField('value', value);
        return;
      }

      const prefix = this.searchInput.substring(0, Math.max(selectionStart, 0));
      const suffix = this.searchInput.substring(Math.max(selectionEnd, 0));

      this.selectionStart = prefix.length;
      this.selectionEnd = this.selectionStart + value.length;
      this.showItems = false;

      this.updateField('value', `${prefix}${value}${suffix}`);
    },

    onSelectionChange() {
      if (!this.$el.contains(document.activeElement) && this.$refs.combobox) {
        return;
      }

      if (!this.searchInput) {
        this.showItems = true;
        return;
      }

      const { selectionStart, selectionEnd } = this.$refs.combobox.$refs.input;

      this.selectionStart = selectionStart;
      this.selectionEnd = selectionEnd;

      const variableGroup = matchPayloadVariableBySelection(this.searchInput, selectionStart, selectionEnd);

      if (!variableGroup) {
        this.showItems = this.searchInput[selectionStart - 1] === '{';

        if (this.showItems) {
          this.selectionStart = selectionStart - 1;
          this.selectionEnd = selectionStart;
        }

        this.activeValue = undefined;
        return;
      }

      const [value] = variableGroup;
      this.activeValue = value;
      this.showItems = true;

      this.selectionStart = variableGroup.index;
      this.selectionEnd = this.selectionStart + value.length;
    },

    removeCondition(key) {
      this.$emit('remove', key);
    },
  },
};
</script>
