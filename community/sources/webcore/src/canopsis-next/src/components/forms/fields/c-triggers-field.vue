<template lang="pug">
  v-select(
    v-validate="'required'",
    v-field="value",
    :items="availableTriggers",
    :disabled="disabled",
    :label="label || $tc('common.trigger', 2)",
    :error-messages="errorMessages",
    :name="name",
    item-disabled="deprecated",
    multiple,
    chips
  )
    template(#selection="{ item, index }")
      v-tooltip(:disabled="!item.deprecated", top)
        template(#activator="{ on }")
          v-chip(
            v-on="on",
            :class="{ 'error--text': item.deprecated }",
            :close="item.deprecated",
            @input="removeItemFromArray(index)"
          ) {{ item.text }}
        span {{ $t('common.deprecatedTrigger') }}
    template(#item="{ item, tile, parent }")
      v-list-tile(v-bind="tile.props", v-on="tile.on")
        v-list-tile-action
          v-checkbox(:input-value="tile.props.value", :color="parent.color")
        v-list-tile-content {{ item.text }}
        v-list-tile-action(v-if="item.helpText")
          c-help-icon(:text="item.helpText", color="info", size="20", top)
</template>

<script>
import { TRIGGERS, PRO_TRIGGERS } from '@/constants';

import { isDeprecatedTrigger } from '@/helpers/entities/scenarios';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { formArrayMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin, entitiesInfoMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'triggers',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    triggers: {
      type: Array,
      default: () => Object.values(TRIGGERS),
    },
  },
  computed: {
    availableTriggers() {
      return Object.values(this.triggers)
        .reduce((acc, type) => {
          if (!PRO_TRIGGERS.includes(type) || this.isProVersion) {
            const { text, helpText } = this.$t(`common.triggers.${type}`);

            acc.push({
              text,
              helpText,
              value: type,
              deprecated: isDeprecatedTrigger(type),
            });
          }

          return acc;
        }, []);
    },

    deprecatedValues() {
      return this.value.filter(isDeprecatedTrigger);
    },

    errorMessages() {
      return this.errors.collect(this.name, null, false)
        .map((item) => {
          const messageMap = {
            max_value: this.$tc(
              'errors.triggerMustNotUsed',
              this.deprecatedValues.length,
              { field: this.deprecatedValues.join(', ') },
            ),
          };

          return messageMap[item.rule] ?? item.msg;
        });
    },
  },
  created() {
    this.attachMaxValueRule();
  },
  beforeDestroy() {
    this.detachRules();
  },
  methods: {
    attachMaxValueRule() {
      this.$validator.attach({
        name: this.name,
        rules: 'max_value:0',
        getter: () => this.deprecatedValues.length,
        vm: this,
      });
    },

    detachRules() {
      this.$validator.detach(this.name);
    },
  },
};
</script>
