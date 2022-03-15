<template lang="pug">
  div
    v-tabs.filter-editor(v-model="activeTab", slider-color="primary", centered)
      v-tab(:disabled="advancedJsonWasChanged || errors.has('advancedJson')") {{ $t('filterEditor.tabs.visualEditor') }}
      v-tab-item
        v-container.pa-1
          div.position-relative
            v-fade-transition(v-if="filterHintsPending || entityCategoriesPending", key="progress", mode="out-in")
              v-progress-linear.progress-linear-absolute--top(height="2", indeterminate)
            filter-group(
              v-field="form",
              :possible-fields="possibleFields",
              is-initial,
              @input="resetFilterValidator"
            )
      v-tab(@click="openAdvancedTab") {{ $t('filterEditor.tabs.advancedEditor') }}
      v-tab-item
        c-json-field(
          :value="advancedJson",
          :label="$t('filterEditor.tabs.advancedEditor')",
          name="advancedJson",
          rows="10",
          validate-on="button",
          @input="updateJson"
        )
    v-alert(:value="errors.has('filter')", type="error") {{ $t('filterEditor.errors.required') }}
</template>

<script>
import { get } from 'lodash';

import { ENTITIES_TYPES, FILTER_OPERATORS, PATTERN_INPUT_TYPES, ENTITY_TYPES, MAX_LIMIT } from '@/constants';

import { filterToForm, formToFilter } from '@/helpers/forms/filter';
import { checkIfGroupIsEmpty } from '@/helpers/filter/editor/filter-check';

import { entitiesFilterHintsMixin } from '@/mixins/entities/associative-table/filter-hints';
import entitiesEntityCategoryMixin from '@/mixins/entities/entity-category';
import { formValidationHeaderMixin } from '@/mixins/form';

import FilterGroup from './partial/filter-group.vue';

/**
 * Component to create new MongoDB filter
 */
export default {
  inject: ['$validator'],
  components: {
    FilterGroup,
  },
  mixins: [
    entitiesFilterHintsMixin,
    entitiesEntityCategoryMixin,
    formValidationHeaderMixin,
  ],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    entitiesType: {
      type: String,
      default: ENTITIES_TYPES.alarm,
      validator: value => [
        ENTITIES_TYPES.alarm,
        ENTITIES_TYPES.entity,
        ENTITIES_TYPES.pbehavior,
      ].includes(value),
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      activeTab: 0,
      advancedJson: '{}',
      filterHints: {},
      filterHintsPending: false,
    };
  },
  computed: {
    advancedJsonWasChanged() {
      return get(this.fields, ['advancedJson', 'changed']);
    },

    categoryHint() {
      return {
        name: this.$t('common.category'),
        value: 'category',
        operatorProps: {
          items: [
            FILTER_OPERATORS.equal,
            FILTER_OPERATORS.notEqual,
            FILTER_OPERATORS.hasOneOf,
            FILTER_OPERATORS.hasNot,
          ],
        },
        valueProps: {
          types: [{ value: PATTERN_INPUT_TYPES.string }],
          items: this.entityCategories,
          itemText: 'name',
          itemValue: '_id',
        },
      };
    },

    impactStateHint() {
      return {
        name: this.$t('common.priority'),
        value: 'impact_state',
        operatorProps: {
          items: [
            FILTER_OPERATORS.greater,
            FILTER_OPERATORS.less,
            FILTER_OPERATORS.equal,
            FILTER_OPERATORS.notEqual,
          ],
        },
        valueProps: {
          types: [{ value: PATTERN_INPUT_TYPES.number }],
        },
      };
    },

    impactLevelHint() {
      return {
        name: this.$t('common.impactLevel'),
        value: 'impact_level',
        operatorProps: {
          items: [
            FILTER_OPERATORS.greater,
            FILTER_OPERATORS.less,
            FILTER_OPERATORS.equal,
            FILTER_OPERATORS.notEqual,
          ],
        },
        valueProps: {
          types: [{ value: PATTERN_INPUT_TYPES.number, defaultValue: 1 }],
        },
        valueValidationRules: { required: true, between: [1, 10] },
      };
    },

    typeHint() {
      return {
        name: this.$t('common.type'),
        value: 'type',
        operatorProps: {
          items: [FILTER_OPERATORS.equal, FILTER_OPERATORS.notEqual],
        },
        valueProps: {
          types: [{ value: PATTERN_INPUT_TYPES.string }],
          items: Object.values(ENTITY_TYPES),
        },
      };
    },

    infosHint() {
      return {
        name: this.$t('common.infos'),
        value: 'infos',
        additionalFieldProps: {
          items: ['description', 'value'],
        },
        valueProps: {
          types: [{ value: PATTERN_INPUT_TYPES.string }],
        },
      };
    },

    defaultAlarmHints() {
      return [
        this.impactStateHint,
        this.categoryHint,
        {
          name: this.$t('filterEditor.hints.alarm.service'),
          value: 'service',
        },
        {
          name: this.$t('filterEditor.hints.alarm.connector'),
          value: 'connector',
        },
        {
          name: this.$t('filterEditor.hints.alarm.connectorName'),
          value: 'connector_name',
        },
        {
          name: this.$t('filterEditor.hints.alarm.component'),
          value: 'component',
        },
        {
          name: this.$t('filterEditor.hints.alarm.resource'),
          value: 'resource',
        },
        this.infosHint,
      ];
    },

    defaultEntityHints() {
      return [
        {
          name: this.$t('common.name'),
          value: 'name',
        },
        this.categoryHint,
        this.impactLevelHint,
        this.typeHint,
        {
          name: this.$t('entity.impact'),
          value: 'impact',
        },
        {
          name: this.$t('entity.depends'),
          value: 'depends',
        },
      ];
    },

    alarmFilterHintsOrDefault() {
      return get(this.filterHints, 'alarm', this.defaultAlarmHints);
    },

    entityFilterHintsOrDefault() {
      return get(this.filterHints, 'entity', this.defaultEntityHints);
    },

    possibleFields() {
      if (this.entitiesType === ENTITIES_TYPES.entity) {
        return this.entityFilterHintsOrDefault;
      }

      return this.alarmFilterHintsOrDefault;
    },

    hasCategoryHint() {
      return this.possibleFields.some(field => field.value === this.categoryHint.value);
    },
  },

  watch: {
    advancedJsonWasChanged(value) {
      if (value) {
        this.resetFilterValidator();
      }
    },

    hasCategoryHint: {
      immediate: true,
      handler(hasCategoryHint) {
        if (hasCategoryHint) {
          this.fetchEntityCategoriesList({ params: { limit: MAX_LIMIT } });
        }
      },
    },
  },
  created() {
    if (this.required && this.$validator) {
      this.$validator.attach({
        name: 'filter',
        rules: 'required:true',
        getter: () => !checkIfGroupIsEmpty(this.form),
        context: () => this,
        vm: this,
      });
    }

    this.fetchList();
  },
  methods: {
    resetFilterValidator() {
      if (this.errors.has('filter')) {
        this.errors.remove('filter');
      }
    },

    openAdvancedTab() {
      if (this.activeTab === 1) {
        return;
      }

      try {
        this.advancedJson = formToFilter(this.form);
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: this.$t('errors.default') });
      }
    },

    updateJson(advancedJson) {
      try {
        this.$emit('input', filterToForm(advancedJson));

        this.advancedJson = advancedJson;
      } catch (err) {
        console.error(err);

        /**
         * We need to use setTimeout instead of $nextTick here because we already used reset inside json-field
         * and $nextTick will not work
         */
        setTimeout(() => {
          this.$validator.flag('advancedJson', {
            touched: true,
          });

          this.errors.add({
            field: 'advancedJson',
            msg: this.$t('filterEditor.errors.cantParseToVisualEditor'),
          });
        }, 0);
      }
    },

    async fetchList() {
      this.filterHintsPending = true;
      this.filterHints = await this.fetchFilterHints();
      this.filterHintsPending = false;
    },
  },
};
</script>

<style lang="scss">
  .filter-editor {
    .v-card {
      box-shadow: 0 0 0 -1px rgba(0, 0, 0, 0.5), 0 1px 5px 0 rgba(0, 0, 0, 0.44), 0 1px 3px 0 rgba(0, 0, 0, 0.42);
    }
  }
</style>
