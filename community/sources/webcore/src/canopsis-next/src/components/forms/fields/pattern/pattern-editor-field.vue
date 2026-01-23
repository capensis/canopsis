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
          class="gap-2"
          align-center
          justify-end
        >
          <pattern-count-message :error="count === 0">
            <span v-html="countMessage" />
          </pattern-count-message>
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
import {
  ref,
  computed,
  watch,
  onMounted,
  onBeforeUnmount,
} from 'vue';
import { isEqual, isEmpty } from 'lodash';

import { PATTERN_CUSTOM_ITEM_VALUE, PATTERN_EDITOR_TABS } from '@/constants';

import { formGroupsToPatternRules, patternsToGroups, patternToForm } from '@/helpers/entities/pattern/form';

import { useModelField } from '@/hooks/form/model-field';
import { useValidator } from '@/hooks/validator/validator';
import { useComponentInstance } from '@/hooks/vue';

import { usePatternCountMessage } from './hooks/pattern-count-message';
import PatternAdvancedEditorField from './pattern-advanced-editor-field.vue';
import PatternGroupsField from './pattern-groups-field.vue';
import PatternCountMessage from './pattern-count-message.vue';

export default {
  inject: ['$validator'],
  components: { PatternCountMessage, PatternGroupsField, PatternAdvancedEditorField },
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
    alarmCounter: {
      type: Object,
      required: false,
    },
    entityCounter: {
      type: Object,
      required: false,
    },
  },
  setup(props, { emit }) {
    const validator = useValidator();
    const { errors } = validator;
    const vm = useComponentInstance();
    const { updateModel, updateField } = useModelField(props, emit);
    const { getCountMessage } = usePatternCountMessage();

    const activeTab = ref(PATTERN_EDITOR_TABS.simple);
    const patternsJson = ref([]);

    const patternGroupsFieldName = computed(() => `${props.name}.groups`);
    const patternJsonFieldName = computed(() => `${props.name}.json`);

    const hasJsonError = computed(() => errors.has(patternJsonFieldName.value));
    const errorMessage = computed(() => errors.collect(props.name)?.join('\n'));

    const isSimpleTab = computed(() => activeTab.value === PATTERN_EDITOR_TABS.simple);
    const isCustomPattern = computed(() => props.patterns.id === PATTERN_CUSTOM_ITEM_VALUE);
    const formDisabled = computed(() => props.disabled || (props.withType && !isCustomPattern.value));

    const checked = computed(() => !isEmpty(props.alarmCounter) || !isEmpty(props.entityCounter));
    const count = computed(() => props.alarmCounter?.count ?? props.entityCounter?.count ?? 0);
    const overLimit = computed(() => props.alarmCounter?.over_limit || props.entityCounter?.over_limit || false);
    const countMessage = computed(() => getCountMessage(props.alarmCounter, props.entityCounter));

    /**
     * Watches for changes in pattern groups and removes validation errors when groups change.
     *
     * @param {Array} groups - New pattern groups value.
     * @param {Array} oldGroups - Previous pattern groups value.
     */
    watch(
      () => props.patterns.groups,
      (groups, oldGroups) => {
        if (!isEqual(groups, oldGroups)) {
          errors.remove(props.name);
        }
      },
    );

    watch(activeTab, (newTab) => {
      if (newTab === PATTERN_EDITOR_TABS.advanced) {
        patternsJson.value = formGroupsToPatternRules(props.patterns.groups);
      }
    });

    /**
     * Updates the pattern with new pattern data from the pattern field.
     *
     * @param {Object} pattern - Pattern object containing pattern data.
     */
    const updatePattern = (pattern) => {
      const { groups } = patternToForm(pattern);

      updateModel({
        ...props.patterns,
        is_corporate: pattern.is_corporate,
        id: pattern._id,
        groups,
      });
    };

    /**
     * Updates the pattern to use custom pattern type.
     */
    const updatePatternToCustom = () => updateField('id', PATTERN_CUSTOM_ITEM_VALUE);

    /**
     * Updates pattern groups from advanced editor patterns and syncs patterns JSON.
     *
     * @param {Array} patterns - Array of pattern rules from advanced editor.
     */
    const updateGroupsFromPatterns = (patterns) => {
      updateField('groups', patternsToGroups(patterns));

      patternsJson.value = patterns;
    };

    onMounted(() => validator.attach({
      name: props.name,
      getter: () => props.patterns.length,
      vm,
    }));

    onBeforeUnmount(() => validator.detach(props.name));

    return {
      activeTab,
      patternsJson,
      patternGroupsFieldName,
      patternJsonFieldName,
      hasJsonError,
      errorMessage,
      isSimpleTab,
      isCustomPattern,
      formDisabled,
      checked,
      count,
      overLimit,
      countMessage,
      updatePattern,
      updatePatternToCustom,
      updateGroupsFromPatterns,
    };
  },
};
</script>
