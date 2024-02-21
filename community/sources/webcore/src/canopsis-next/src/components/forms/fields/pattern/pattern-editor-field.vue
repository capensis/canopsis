<template>
  <v-layout column>
    <c-pattern-field
      v-if="withType"
      :value="patterns.id"
      :type="type"
      :disabled="disabled || readonly"
      class="mb-2"
      return-object
      required
      @input="updatePattern"
    />
    <v-tabs
      v-if="!withType || patterns.id"
      v-model="activeTab"
      slider-color="primary"
      centered
    >
      <v-tab
        :disabled="!isSimpleTab && hasJsonError"
        :href="`#${$constants.PATTERN_EDITOR_TABS.simple}`"
      >
        {{ $t('pattern.simpleEditor') }}
      </v-tab>
      <v-tab-item :value="$constants.PATTERN_EDITOR_TABS.simple">
        <pattern-groups-field
          v-field="patterns.groups"
          :disabled="formDisabled"
          :readonly="readonly"
          :name="patternGroupsFieldName"
          :type="type"
          :required="required"
          :attributes="attributes"
          class="mt-2"
        />
      </v-tab-item>
      <v-tab :href="`#${$constants.PATTERN_EDITOR_TABS.advanced}`">
        {{ $t('pattern.advancedEditor') }}
      </v-tab>
      <v-tab-item :value="$constants.PATTERN_EDITOR_TABS.advanced">
        <pattern-advanced-editor-field
          :value="patternsJson"
          :disabled="readonly || disabled || !isCustomPattern"
          :attributes="attributes"
          :name="patternJsonFieldName"
          @input="updateGroupsFromPatterns"
        />
      </v-tab-item>
    </v-tabs>
    <template v-if="!readonly">
      <v-layout
        align-center
        justify-end
      >
        <v-btn
          v-if="withType && !isCustomPattern"
          :disabled="disabled"
          class="mr-0"
          color="primary"
          @click="updatePatternToCustom"
        >
          {{ $t('common.edit') }}
        </v-btn>
        <v-layout
          v-if="checked"
          align-center
          justify-end
        >
          <pattern-count-message
            :error="count === 0"
            :message="$tc('common.itemFound', count, { count })"
          />
          <slot name="append-count" />
        </v-layout>
      </v-layout>
      <v-flex class="mt-2">
        <v-alert
          :value="!!errorMessage"
          class="pre-wrap"
          type="error"
          transition="fade-transition"
        >
          {{ errorMessage }}
        </v-alert>
        <v-alert
          :value="overLimit"
          type="warning"
          transition="fade-transition"
        >
          <span>{{ $t('pattern.errors.countOverLimit', { count }) }}</span>
        </v-alert>
      </v-flex>
    </template>
  </v-layout>
</template>

<script>
import { isEqual, isEmpty } from 'lodash';

import { PATTERN_CUSTOM_ITEM_VALUE, PATTERN_EDITOR_TABS } from '@/constants';

import { formGroupsToPatternRules, patternsToGroups, patternToForm } from '@/helpers/entities/pattern/form';

import { formMixin, validationChildrenMixin } from '@/mixins/form';

import PatternAdvancedEditorField from './pattern-advanced-editor-field.vue';
import PatternGroupsField from './pattern-groups-field.vue';
import PatternCountMessage from './pattern-count-message.vue';

export default {
  inject: ['$validator'],
  components: { PatternCountMessage, PatternGroupsField, PatternAdvancedEditorField },
  mixins: [formMixin, validationChildrenMixin],
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
      required: true,
    },
    attributes: {
      type: Array,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'patterns',
    },
    checkCountName: {
      type: String,
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    type: {
      type: String,
      required: false,
    },
    withType: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    counter: {
      type: Object,
      required: false,
    },
  },
  data() {
    return {
      activeTab: PATTERN_EDITOR_TABS.simple,
      patternsJson: [],
    };
  },
  computed: {
    patternGroupsFieldName() {
      return `${this.name}.groups`;
    },

    patternJsonFieldName() {
      return `${this.name}.json`;
    },

    hasJsonError() {
      return this.errors.has(this.patternJsonFieldName);
    },

    errorMessage() {
      return this.errors.collect(this.name)?.join('\n');
    },

    isSimpleTab() {
      return this.activeTab === PATTERN_EDITOR_TABS.simple;
    },

    formDisabled() {
      return this.disabled || (this.withType && !this.isCustomPattern);
    },

    isCustomPattern() {
      return this.patterns.id === PATTERN_CUSTOM_ITEM_VALUE;
    },

    checked() {
      return !isEmpty(this.counter);
    },

    count() {
      return this.counter?.count ?? 0;
    },

    overLimit() {
      return this.counter?.over_limit ?? false;
    },
  },
  watch: {
    'patterns.groups': {
      handler(groups, oldGroups) {
        if (!isEqual(groups, oldGroups)) {
          this.errors.remove(this.name);
        }
      },
    },

    activeTab(newTab) {
      if (newTab === PATTERN_EDITOR_TABS.advanced) {
        this.patternsJson = formGroupsToPatternRules(this.patterns.groups);
      }
    },
  },
  created() {
    this.$validator.attach({
      name: this.name,
      getter: () => this.patterns.length,
      vm: this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.name);
  },
  methods: {
    updatePattern(pattern) {
      const { groups } = patternToForm(pattern);

      this.updateModel({
        ...this.patterns,
        is_corporate: pattern.is_corporate,
        id: pattern._id,
        groups,
      });
    },

    updatePatternToCustom() {
      this.updateField('id', PATTERN_CUSTOM_ITEM_VALUE);
    },

    updateGroupsFromPatterns(patterns) {
      this.updateField('groups', patternsToGroups(patterns));

      this.patternsJson = patterns;
    },
  },
};
</script>
