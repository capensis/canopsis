<template>
  <c-layout column>
    <c-pattern-field
      v-if="withType"
      :value="patterns.id"
      :type="type"
      :disabled="disabled || readonly"
      class="mb-2"
      return-object
      required
      @update:value="updatePattern"
    />
    <v-tabs
      v-if="!withType || patterns.id"
      v-model="activeTab"
      align-tabs="center"
      color="primary"
    >
      <v-tab
        :disabled="!isSimpleTab && hasJsonError"
        :value="$constants.PATTERN_EDITOR_TABS.simple"
      >
        {{ $t('pattern.simpleEditor') }}
      </v-tab>

      <v-tab :value="$constants.PATTERN_EDITOR_TABS.advanced">
        {{ $t('pattern.advancedEditor') }}
      </v-tab>
    </v-tabs>

    <v-window
      v-if="!withType || patterns.id"
      v-model="activeTab"
    >
      <v-window-item :value="$constants.PATTERN_EDITOR_TABS.simple">
        <pattern-groups-field
          :groups="patterns.groups"
          :disabled="formDisabled"
          :readonly="readonly"
          :name="patternGroupsFieldName"
          :type="type"
          :required="required"
          :attributes="attributes"
          class="mt-2"
          @update:groups="updateField('groups', $event)"
        />
      </v-window-item>
      <v-window-item :value="$constants.PATTERN_EDITOR_TABS.advanced">
        <pattern-advanced-editor-field
          :value="patternsJson"
          :disabled="readonly || disabled || !isCustomPattern"
          :attributes="attributes"
          :name="patternJsonFieldName"
          @update:value="updateGroupsFromPatterns"
        />
      </v-window-item>
    </v-window>
    <template v-if="!readonly">
      <c-layout
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
        <c-layout
          v-if="checked"
          class="gap-2"
          align-center
          justify-end
        >
          <pattern-count-message
            :alarm-counter="alarmCounter"
            :entity-counter="entityCounter"
          />
          <slot name="append-count" />
        </c-layout>
      </c-layout>
      <div class="mt-2">
        <v-alert
          v-if="errorMessage"
          class="pre-wrap"
          type="error"
        >
          {{ errorMessage }}
        </v-alert>
        <v-alert
          v-if="overLimit"
          type="warning"
        >
          <span>{{ $t('pattern.errors.countOverLimit', { count }) }}</span>
        </v-alert>
      </div>
    </template>
  </c-layout>
</template>

<script>
import {
  ref,
  computed,
  watch,
} from 'vue';
import { useFieldError, useSetFieldError } from 'vee-validate';
import { isEqual, isEmpty } from 'lodash';

import { PATTERN_CUSTOM_ITEM_VALUE, PATTERN_EDITOR_TABS } from '@/constants';

import { formGroupsToPatternRules, patternsToGroups, patternToForm } from '@/helpers/entities/pattern/form';

import { useModelField } from '@/hooks/form/model-field';

import PatternAdvancedEditorField from './pattern-advanced-editor-field.vue';
import PatternGroupsField from './pattern-groups-field.vue';
import PatternCountMessage from './pattern-count-message.vue';

export default {
  components: { PatternCountMessage, PatternGroupsField, PatternAdvancedEditorField },
  emits: ['update:patterns'],
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
    const { updateModel, updateField } = useModelField(props, emit);
    const setFieldError = useSetFieldError();

    const activeTab = ref(PATTERN_EDITOR_TABS.simple);
    const patternsJson = ref([]);

    const patternGroupsFieldName = computed(() => `${props.name}.groups`);
    const patternJsonFieldName = computed(() => `${props.name}.json`);

    const jsonErrorMessage = useFieldError(patternJsonFieldName);
    const errorMessage = useFieldError(computed(() => props.name));
    const hasJsonError = computed(() => !!jsonErrorMessage.value);

    const isSimpleTab = computed(() => activeTab.value === PATTERN_EDITOR_TABS.simple);
    const isCustomPattern = computed(() => props.patterns.id === PATTERN_CUSTOM_ITEM_VALUE);
    const formDisabled = computed(() => props.disabled || (props.withType && !isCustomPattern.value));

    const checked = computed(() => !isEmpty(props.alarmCounter) || !isEmpty(props.entityCounter));
    const count = computed(() => props.alarmCounter?.count ?? props.entityCounter?.count ?? 0);
    const overLimit = computed(() => props.alarmCounter?.over_limit || props.entityCounter?.over_limit || false);

    /**
     * Watches for changes in pattern groups and removes validation errors when groups change.
     *
     * @param {Array} groups - New pattern groups value.
     * @param {Array} oldGroups - Previous pattern groups value.
     */
    watch(
      () => props.patterns.groups,
      (groups, oldGroups) => {
        if (hasJsonError.value) {
          setFieldError(patternJsonFieldName.value, undefined);

          if (activeTab.value === PATTERN_EDITOR_TABS.advanced) {
            patternsJson.value = formGroupsToPatternRules(groups);
          }
        }

        if (!isEqual(groups, oldGroups)) {
          setFieldError(props.name, undefined);
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
      updateField,
      updatePattern,
      updatePatternToCustom,
      updateGroupsFromPatterns,
    };
  },
};
</script>
