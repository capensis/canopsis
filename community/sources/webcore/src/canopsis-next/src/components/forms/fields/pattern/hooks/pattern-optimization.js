import { cloneDeep } from 'lodash';
import {
  computed,
  ref,
  unref,
  inject,
  watch,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import {
  ENTITY_PATTERN_FIELDS,
  MODALS,
  PATTERN_OPERATORS,
  PATTERN_OPTIMIZATION_STATUSES,
  PATTERNS_FIELDS,
  PATTERN_TYPES,
  PATTERN_CUSTOM_ITEM_VALUE,
} from '@/constants';

import Observer from '@/services/observer';

import { patternToForm } from '@/helpers/entities/pattern/form';
import { formFilterToPatterns } from '@/helpers/entities/filter/form';
import { isOmitEqual } from '@/helpers/collection';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { usePollingWithPending } from '@/hooks/polling';
import { useValidator } from '@/hooks/validator/validator';
import { useModelField } from '@/hooks/form/model-field';
import { usePatternEntitiesOptimize } from '@/hooks/store/modules/pattern-entities-optimize';

/**
 * Composable hook for managing pattern optimization functionality
 *
 * Handles pattern optimization process including:
 * - Starting and polling optimization requests
 * - Managing optimization state and suggestions
 * - Applying and rejecting optimization suggestions
 * - Showing entity comparison modals
 * - Tracking applied suggestions for form submission
 *
 * @param {Ref|Object} value - Reactive form value containing pattern data
 * @param {Function} emit - Vue emit function for updating parent component
 *
 * @returns {Object} Pattern optimization API
 * @returns {Ref<boolean>} returns.pending - Whether optimization is currently in progress
 * @returns {ComputedRef<string>} returns.failedReason - Reason for optimization failure, if any
 * @returns {ComputedRef<boolean>} returns.successful - Whether optimization completed successfully
 * @returns {ComputedRef<Array>} returns.suggestions - Array of optimization suggestions
 * @returns {ComputedRef<Array>} returns.optimizedFieldsRegexps - Array of optimized field regexps
 * @returns {ComputedRef<boolean>} returns.mayHaveOptimizationSuggestions - Whether pattern may have suggestions
 * @returns {Function} returns.tryOptimization - Starts pattern optimization process
 * @returns {Function} returns.cancelOptimization - Cancels ongoing optimization process
 * @returns {Function} returns.applySuggestion - Applies a specific optimization suggestion
 * @returns {Function} returns.rejectAllSuggestions - Rejects all suggestions and resets optimization state
 * @returns {Function} returns.showEntitiesComparisonModal - Shows modal comparing current pattern with suggestion
 */
export const usePatternOptimization = ({ value, wrapperElement }, emit) => {
  const { t } = useI18n();
  const modals = useModals();
  const popups = usePopups();
  const validator = useValidator();
  const { updateModel } = useModelField({ value }, emit);

  const afterSubmitObserver = inject('$afterSubmitObserver', new Observer());

  const {
    optimizeEntities,
    fetchOptimizeEntitiesStatus,
    updateOptimization,
    removeOptimization,
  } = usePatternEntitiesOptimize();

  const appliedOptimizations = ref([]);

  const optimization = ref(null);
  const optimizationPattern = ref(null);

  const failedReason = computed(() => optimization.value?.failed_reason);
  const successful = computed(() => optimization.value?.status === PATTERN_OPTIMIZATION_STATUSES.success);
  const suggestions = computed(() => (optimization.value?.suggestions ?? []));
  const currentPattern = computed(() => formFilterToPatterns(unref(value), [PATTERNS_FIELDS.entity], false));
  const optimizedFieldsRegexps = computed(() => optimization.value?.optimized_field_regexps ?? []);
  const originalCounter = computed(() => (optimization.value
    ? ({
      ms: optimization.value.original_pattern_ms ?? 0,
      count: optimization.value.original_pattern_count ?? 0,
    })
    : null
  ));

  const originalValue = cloneDeep(unref(value));

  const hasChanges = computed(() => {
    const uwrappedValue = unref(value);

    const patternKeysWithTemplates = Object.keys(uwrappedValue).filter(key => (
      uwrappedValue[key].id !== PATTERN_CUSTOM_ITEM_VALUE
    ));

    return !isOmitEqual(uwrappedValue, originalValue, ['title', ...patternKeysWithTemplates]);
  });

  const hasRegexpInfos = computed(() => {
    const { [PATTERNS_FIELDS.entity]: entityPattern } = unref(value);

    return entityPattern?.groups?.some?.(group => (
      group.rules.some(rule => (
        rule.operator === PATTERN_OPERATORS.regexp
        && rule.attribute === ENTITY_PATTERN_FIELDS.infos
      ))
    ));
  });

  const mayHaveOptimizationSuggestions = computed(() => hasRegexpInfos.value && hasChanges.value);

  const { pending, poll, cancel } = usePollingWithPending({
    startHandler: async () => {
      optimizationPattern.value = cloneDeep(currentPattern.value);

      const response = await optimizeEntities({
        data: currentPattern.value,
      });

      optimization.value = response;

      return response;
    },
    processHandler: async ({ _id: id }, resolve) => {
      const response = await fetchOptimizeEntitiesStatus({
        id,
      });

      optimization.value = response;

      if (
        [
          PATTERN_OPTIMIZATION_STATUSES.success,
          PATTERN_OPTIMIZATION_STATUSES.failed,
        ].includes(response.status)
      ) {
        return resolve(response);
      }

      return response;
    },
  });

  /**
   * Rejects all pattern optimization suggestions and resets optimization state
   *
   * Clears the optimization object, effectively dismissing all suggestions.
   */
  const rejectAllSuggestions = () => optimization.value = null;

  /**
   * Starts pattern optimization process
   *
   * Validates form fields (excluding title) and initiates optimization polling.
   * Sets optimization state to failed if validation fails or an error occurs.
   */
  const tryOptimization = async () => {
    try {
      rejectAllSuggestions();

      const unwrappedWrapperElement = unref(wrapperElement);

      const fieldsWithoutTitle = validator.fields.items.filter(({ name, el }) => (
        name !== 'title' && unwrappedWrapperElement.contains(el)
      ));

      const isValid = await validator.validateAll(fieldsWithoutTitle);

      if (!isValid) {
        return;
      }

      optimization.value = await poll();
    } catch (err) {
      console.error(err);

      optimization.value = {
        status: PATTERN_OPTIMIZATION_STATUSES.failed,
        failed_reason: err.message,
      };
    }
  };

  /**
   * Cancels ongoing pattern optimization process
   *
   * Stops polling, removes optimization request from server if exists,
   * and resets optimization state and applied suggestion index.
   */
  const cancelOptimization = () => {
    cancel();

    if (optimization.value?._id) {
      removeOptimization({ id: optimization.value._id });
    }

    rejectAllSuggestions();

    popups.info({
      text: t('pattern.optimizationCancelled'),
      autoClose: 5000,
    });
  };

  /**
   * Applies a pattern optimization suggestion
   *
   * Converts suggestion pattern to form format, updates the form model,
   * tracks the applied suggestion index, and rejects all suggestions.
   *
   * @param {number} index - The index of the suggestion in the suggestions array
   */
  const applySuggestion = (index) => {
    const entityPattern = patternToForm({
      type: PATTERN_TYPES.entity,
      [PATTERNS_FIELDS.entity]: suggestions.value[index][PATTERNS_FIELDS.entity],
    });

    appliedOptimizations.value.push({
      index,
      id: optimization.value?._id,
    });

    updateModel({
      ...unref(value),
      [PATTERNS_FIELDS.entity]: entityPattern,
    });

    rejectAllSuggestions();
  };

  /**
   * Shows modal with entities comparison between current pattern and suggestion
   *
   * Opens entities comparison modal displaying differences between
   * current pattern and the provided suggestion pattern.
   *
   * @param {Object} suggestion - The suggestion object containing pattern groups
   */
  const showEntitiesComparisonModal = suggestion => modals.show({
    name: MODALS.entitiesComparison,
    config: {
      currentPattern: currentPattern.value[PATTERNS_FIELDS.entity],
      suggestionPattern: suggestion?.entity_pattern ?? [],
    },
  });

  /**
   * Handles form submission callback
   *
   * Updates optimization request on server to mark applied suggestion as accepted
   * if a suggestion was applied and optimization request exists.
   */
  const afterSubmit = () => {
    if (appliedOptimizations.value?.length) {
      appliedOptimizations.value.forEach(({ id, index }) => (
        updateOptimization({
          id,
          data: {
            accept: true,
            index,
          },
        })
      ));
    }
  };

  watch(value, rejectAllSuggestions);

  onMounted(() => afterSubmitObserver?.register?.(afterSubmit));
  onBeforeUnmount(() => afterSubmitObserver?.unregister?.(afterSubmit));

  return {
    pending,
    failedReason,
    successful,
    suggestions,
    optimizedFieldsRegexps,
    originalCounter,
    mayHaveOptimizationSuggestions,

    tryOptimization,
    cancelOptimization,
    applySuggestion,
    rejectAllSuggestions,
    showEntitiesComparisonModal,
  };
};
