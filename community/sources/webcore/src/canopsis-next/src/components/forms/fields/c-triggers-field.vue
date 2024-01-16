<template>
  <v-select
    v-validate="'required'"
    :value="value"
    :items="preparedTriggers"
    :disabled="disabled"
    :label="label || $tc('common.trigger', 2)"
    :error-messages="errorMessages"
    :name="name"
    item-value="type"
    item-disabled="deprecated"
    multiple
    chips
    return-object
    @change="changeValue"
  >
    <template #selection="{ item, index }">
      <v-tooltip
        :disabled="!item.deprecated"
        top
      >
        <template #activator="{ on }">
          <v-chip
            v-on="on"
            :class="getSelectedClass(item)"
            :close="item.deprecated"
            @click:close="removeItemFromArray(index)"
          >
            {{ getSelectedText(item) }}
          </v-chip>
        </template>
        <span>{{ $t('common.deprecatedTrigger') }}</span>
      </v-tooltip>
    </template>
    <template #item="{ item, attrs, on, parent }">
      <v-list-item
        v-bind="attrs"
        :active-class="errors.has(getAdditionalValueFieldName(item.type)) ? 'error--text' : tile.props.activeClass"
        @click="on.click"
      >
        <v-list-item-action>
          <v-checkbox
            :input-value="attrs.value"
            :color="parent.color"
            hide-details
          />
        </v-list-item-action>
        <v-list-item-content>
          <v-layout
            class="fill-width"
            align-center
            justify-space-between
          >
            <v-flex>{{ item.text }}</v-flex>
            <component
              class="ml-1"
              v-if="additionalValuesComponentsByTypes[item.type]"
              v-bind="additionalValuesComponentsByTypes[item.type].bind"
              v-on="additionalValuesComponentsByTypes[item.type].on"
              :is="additionalValuesComponentsByTypes[item.type].is"
              :disabled="!attrs.value"
              @click.prevent.stop=""
            />
          </v-layout>
        </v-list-item-content>
        <v-list-item-action v-if="item.helpText">
          <c-help-icon
            :text="item.helpText"
            color="info"
            size="20"
            top
          />
        </v-list-item-action>
      </v-list-item>
    </template>
  </v-select>
</template>

<script>
import { find } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { TRIGGERS_TYPES, PRO_TRIGGERS } from '@/constants';

import { setSeveralFields } from '@/helpers/immutable';
import { isDeprecatedTrigger } from '@/helpers/entities/scenario/form';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { formArrayMixin } from '@/mixins/form';

const { mapGetters } = createNamespacedHelpers('info');

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
    types: {
      type: Array,
      default: () => Object.values(TRIGGERS_TYPES),
    },
  },
  data() {
    return {
      preparedTriggersByTypes: {},
    };
  },
  computed: {
    ...mapGetters({
      eventsCountTriggerDefaultThreshold: 'eventsCountTriggerDefaultThreshold',
    }),

    preparedTriggers() {
      return this.types.reduce((acc, type) => {
        if (!PRO_TRIGGERS.includes(type) || this.isProVersion) {
          const { text, helpText } = this.$t(`common.triggers.${type}`);

          acc.push({
            ...this.preparedTriggersByTypes[type],

            text,
            helpText,
            deprecated: isDeprecatedTrigger(type),
          });
        }

        return acc;
      }, []);
    },

    deprecatedValues() {
      return this.value.filter(({ type }) => isDeprecatedTrigger(type));
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

    defaultAdditionalValuesByTriggers() {
      return {
        [TRIGGERS_TYPES.eventscount]: this.eventsCountTriggerDefaultThreshold ?? '',
      };
    },

    additionalValuesKeysByTriggers() {
      return {
        [TRIGGERS_TYPES.eventscount]: 'threshold',
      };
    },

    additionalValuesComponentsByTypes() {
      const additionalValuesComponents = [
        {
          type: TRIGGERS_TYPES.eventscount,
          is: 'c-number-field',
          bind: {
            class: ['mt-0', 'pt-2'],
            hideDetails: true,
            min: 1,
            required: true,
          },
        },
      ];

      return additionalValuesComponents.reduce((acc, { type, ...rest }) => {
        const trigger = this.preparedTriggersByTypes[type];
        const name = this.getAdditionalValueFieldName(type);

        if (!trigger) {
          return acc;
        }

        acc[type] = setSeveralFields(rest, {
          'bind.name': name,
          'bind.value': trigger[this.additionalValuesKeysByTriggers[type]],
          'bind.label': this.$t(`common.triggers.${type}.additionalFieldLabel`),
          'bind.errorMessages': this.errors.collect(name),
          'on.input': () => value => this.changeAdditionalValue(type, value),
        });

        return acc;
      }, {});
    },
  },
  watch: {
    types(types) {
      this.setPreparedTriggersByTypes(types);
    },
  },
  created() {
    this.attachMaxValueRule();
    this.setPreparedTriggersByTypes(this.types);
  },
  beforeDestroy() {
    this.detachRules();
  },
  methods: {
    /**
     * Get name property for additional value field
     *
     * @param {string} type
     * @return {string}
     */
    getAdditionalValueFieldName(type) {
      return `${this.name}.${type}.additionalValue`;
    },

    /**
     * Get CSS class for selected item
     *
     * @param {boolean} deprecated
     * @param {string} type
     * @return {{ error: boolean, 'error--text': boolean }}
     */
    getSelectedClass({ deprecated, type }) {
      return {
        'error--text': deprecated,
        error: this.errors.has(this.getAdditionalValueFieldName(type)),
      };
    },

    /**
     * Get text for selected item
     *
     * @param {string} type
     * @param {string} text
     * @param {number | string | boolean} [additionalValue]
     * @return {VueI18n.TranslateResult|*}
     */
    getSelectedText({ type, text, [this.additionalValuesKeysByTriggers[type]]: additionalValue } = {}) {
      const messageKey = `common.triggers.${type}.selectedText`;

      return this.$te(messageKey)
        ? this.$t(messageKey, { additionalValue })
        : text;
    },

    /**
     * Set preparedTriggersByTypes by types array
     *
     * @param {string[]} [types = []]
     */
    setPreparedTriggersByTypes(types = []) {
      this.preparedTriggersByTypes = types.reduce((acc, type) => {
        const additionalValueKey = this.additionalValuesKeysByTriggers[type];
        const {
          [additionalValueKey]: additionalValue = this.defaultAdditionalValuesByTriggers[type],
        } = find(this.value, { type }) ?? {};

        acc[type] = {
          type,
          [additionalValueKey]: additionalValue,
        };

        return acc;
      }, {});
    },

    /**
     * Change value on select field handler
     *
     * @param {Trigger[]} value
     */
    changeValue(value = []) {
      this.updateModel(value.map(({ type }) => this.preparedTriggersByTypes[type]));
    },

    /**
     * Change additional value handler
     *
     * @param {string} type
     * @param {number | string | boolean} additionalValue
     */
    changeAdditionalValue(type, additionalValue) {
      this.$set(
        this.preparedTriggersByTypes[type],
        this.additionalValuesKeysByTriggers[type],
        additionalValue,
      );

      this.updateModel(this.value.map(trigger => (
        trigger.type === type
          ? this.preparedTriggersByTypes[type]
          : trigger
      )));
    },

    /**
     * Attach rule for deprecatedValues checking into validator
     */
    attachMaxValueRule() {
      this.$validator.attach({
        name: this.name,
        rules: 'max_value:0',
        getter: () => this.deprecatedValues.length,
        vm: this,
      });
    },

    /**
     * Detach rule for deprecatedValues checking from validator
     */
    detachRules() {
      this.$validator.detach(this.name);
    },
  },
};
</script>
