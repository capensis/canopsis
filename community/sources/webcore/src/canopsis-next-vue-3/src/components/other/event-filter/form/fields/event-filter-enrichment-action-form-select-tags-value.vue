<template>
  <c-select
    :model-value="value"
    :label="$t('common.value')"
    :items="items"
    :name="name"
    class="ml-2"
    required
    @update:model-value="updateValue"
  />
</template>

<script>
import { find, isEqual } from 'lodash';

import { EVENT_FILTER_SET_TAGS_REGEX } from '@/constants';
import { formMixin } from '@/mixins/form';
import { validatorMixin } from '@/mixins/validator';

export default {
  emits: ['update:value'],
  mixins: [formMixin, validatorMixin],
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

    // TODO: Remove this computed property when vee-validate will be updated to 4.10.0
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
          this.value &&
          !isEqual(items, oldItems) &&
          items.every(({ value }) => value !== this.value)
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
  beforeUnmount() {
    this.detachRegexRule();
  },
  methods: {
    updateValue(value) {
      this.updateModel(value);
      this.validateRegex();
    },

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
