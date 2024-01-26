<template>
  <v-select
    class="ml-2"
    v-field="value"
    v-validate="'required'"
    :label="$t('common.value')"
    :error-messages="errorMessages"
    :items="items"
    :name="name"
    @change="validateRegex"
  />
</template>

<script>
import { find, isEqual } from 'lodash';

import { EVENT_FILTER_SET_TAGS_REGEX } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'tags_value',
    },
  },
  computed: {
    selectedItem() {
      return find(this.items, { value: this.value });
    },

    regexRuleFieldName() {
      return `${this.name}.regex`;
    },

    errorMessages() {
      const messages = this.errors.collect(this.name);

      if (this.errors.has(this.regexRuleFieldName)) {
        messages.push(this.$t('eventFilter.validation.incorrectRegexOnSetTagsValue'));
      }

      return messages;
    },
  },
  watch: {
    items: {
      immediate: true,
      handler(items, oldItems) {
        if (
          this.value
          && !isEqual(items, oldItems)
          && items.every(({ value }) => value !== this.value)
        ) {
          this.updateField('value', '');
        }

        this.validateRegex();
      },
    },
  },
  mounted() {
    this.attachRegexRule();
  },
  beforeDestroy() {
    this.detachRegexRule();
  },
  methods: {
    attachRegexRule() {
      this.$validator.attach({
        name: this.regexRuleFieldName,
        rules: {
          regex: EVENT_FILTER_SET_TAGS_REGEX,
        },
        getter: () => this.selectedItem?.valueForValidation ?? '',
        vm: this,
      });
    },

    detachRegexRule() {
      this.$validator.detach(this.regexRuleFieldName);
    },

    validateRegex() {
      if (!this.$validator.fields.find({ name: this.regexRuleFieldName })) {
        return;
      }

      this.$validator.validate(this.regexRuleFieldName);
    },
  },
};
</script>
