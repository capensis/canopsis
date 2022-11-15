<template lang="pug">
  v-layout(column)
    c-pattern-field.mb-2(
      v-if="!patterns.old_mongo_query && withType",
      :value="patterns.id",
      :type="type",
      :disabled="disabled || readonly",
      return-object,
      required,
      @input="updatePattern"
    )

    v-tabs(
      v-if="!withType || patterns.id",
      v-model="activeTab",
      slider-color="primary",
      centered
    )
      v-tab(
        :disabled="!isSimpleTab && hasJsonError",
        :href="`#${$constants.PATTERN_EDITOR_TABS.simple}`"
      ) {{ $t('pattern.simpleEditor') }}
      v-tab-item(:value="$constants.PATTERN_EDITOR_TABS.simple")
        v-layout(v-if="patterns.old_mongo_query", justify-center, wrap)
          v-flex.pt-2(xs8)
            div.error--text.text-xs-center {{ $t('pattern.errors.oldPattern') }}
          v-flex.pt-2(v-if="!readonly && !disabled", xs12)
            v-layout(justify-center)
              v-btn(color="primary", @click="discardPattern") {{ $t('pattern.discard') }}
        c-pattern-groups-field.mt-2(
          v-else,
          v-field="patterns.groups",
          :disabled="formDisabled",
          :readonly="readonly",
          :name="patternGroupsFieldName",
          :type="type",
          :required="required",
          :attributes="attributes"
        )

      v-tab(:href="`#${$constants.PATTERN_EDITOR_TABS.advanced}`") {{ $t('pattern.advancedEditor') }}
      v-tab-item(:value="$constants.PATTERN_EDITOR_TABS.advanced", lazy)
        c-json-field(
          v-if="patterns.old_mongo_query",
          :value="patterns.old_mongo_query",
          :label="$t('pattern.advancedEditor')",
          readonly,
          rows="10"
        )
        c-pattern-advanced-editor-field(
          v-else,
          :value="patternsJson",
          :disabled="readonly || disabled || !isCustomPattern",
          :attributes="attributes",
          :name="patternJsonFieldName",
          @input="updateGroupsFromPatterns"
        )

    template(v-if="!readonly && !patterns.old_mongo_query")
      v-layout(align-center, justify-end)
        div(v-if="checkCountName")
          span.mr-2(
            v-show="patternsChecked",
            :class="{ 'error--text': itemsCount === 0 }"
          ) {{ $tc('common.itemFound', itemsCount, { count: itemsCount }) }}
          v-btn.primary.mx-0(
            :disabled="disabled || patternsChecked || hasChildrenError || !patterns.groups.length",
            :loading="checkPatternsPending",
            @click="checkPatterns"
          ) {{ $t('common.checkPattern') }}
        v-btn.mr-0(
          v-if="withType && !isCustomPattern",
          :disabled="disabled",
          color="primary",
          @click="updatePatternToCustom"
        ) {{ $t('common.edit') }}

      v-flex
        v-alert.pre-wrap(v-if="errorMessage", value="true") {{ errorMessage }}
        v-alert(
          v-if="patternsChecked && itemsOverLimit",
          value="true",
          type="warning",
          transition="fade-transition"
        )
          span {{ $t('pattern.errors.countOverLimit', { count: itemsCount }) }}
</template>

<script>
import { isEqual } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { PATTERN_CUSTOM_ITEM_VALUE, PATTERN_EDITOR_TABS } from '@/constants';

import { formGroupsToPatternRules, patternsToGroups, patternToForm } from '@/helpers/forms/pattern';

import { formMixin, validationChildrenMixin } from '@/mixins/form';

const { mapActions } = createNamespacedHelpers('pattern');

export default {
  inject: ['$validator'],
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
  },
  data() {
    return {
      checkPatternsPending: false,
      patternsChecked: false,
      itemsCount: undefined,
      itemsOverLimit: false,
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
  },
  watch: {
    'patterns.groups': {
      handler(groups, oldGroups) {
        if (!isEqual(groups, oldGroups)) {
          this.patternsChecked = false;

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
    ...mapActions({
      checkPatternsCount: 'checkPatternsCount',
    }),

    discardPattern() {
      this.updateModel(patternToForm({ id: PATTERN_CUSTOM_ITEM_VALUE }));
    },

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

    async checkPatterns() {
      try {
        this.checkPatternsPending = true;

        const isFormValid = await this.validateChildren();

        if (isFormValid) {
          const { [this.checkCountName]: patternsCount } = await this.checkPatternsCount({
            data: {
              [this.checkCountName]: formGroupsToPatternRules(this.patterns.groups),
            },
          });

          const { count, over_limit: overLimit } = patternsCount;

          this.itemsCount = count;
          this.itemsOverLimit = overLimit;
          this.patternsChecked = true;
        }
      } catch (err) {
        const { [this.checkCountName]: error } = err || {};

        if (error) {
          this.errors.add({ field: this.name, msg: error });
        }

        this.patternsChecked = false;
      } finally {
        this.checkPatternsPending = false;
      }
    },
  },
};
</script>
