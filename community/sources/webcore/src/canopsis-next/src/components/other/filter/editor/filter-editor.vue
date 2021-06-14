<template lang="pug">
  div(data-test="filterEditor")
    v-tabs.filter-editor(v-model="activeTab", slider-color="primary", centered)
      v-tab(:disabled="advancedJsonWasChanged || errors.has('advancedJson')") {{ $t('filterEditor.tabs.visualEditor') }}
      v-tab-item
        v-container.pa-1
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

import { ENTITIES_TYPES } from '@/constants';

import { filterToForm, formToFilter } from '@/helpers/forms/filter';
import { checkIfGroupIsEmpty } from '@/helpers/filter/editor/filter-check';


import filterHintsMixin from '@/mixins/entities/filter-hint';
import formValidationHeaderMixin from '@/mixins/form/validation-header';

import FilterGroup from './partial/filter-group.vue';

/**
 * Component to create new MongoDB filter
 */
export default {
  inject: ['$validator'],
  components: {
    FilterGroup,
  },
  mixins: [filterHintsMixin, formValidationHeaderMixin],
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
    };
  },
  computed: {
    advancedJsonWasChanged() {
      return get(this.fields, ['advancedJson', 'changed']);
    },

    defaultAlarmHints() {
      return [
        {
          name: 'Connector',
          value: 'connector',
        },
        {
          name: 'Connector name',
          value: 'connector_name',
        },
        {
          name: 'Component',
          value: 'component',
        },
        {
          name: 'Resource',
          value: 'resource',
        },
      ];
    },

    defaultEntityHints() {
      return [
        {
          name: 'Name',
          value: 'name',
        },
        {
          name: 'Type',
          value: 'type',
        },
        {
          name: 'Impact',
          value: 'impact',
        },
        {
          name: 'Depends',
          value: 'depends',
        },
      ];
    },

    alarmFilterHintsOrDefault() {
      return this.alarmFilterHints || this.defaultAlarmHints;
    },

    entityFilterHintsOrDefault() {
      return this.entityFilterHints || this.defaultEntityHints;
    },

    possibleFields() {
      if (this.entitiesType === ENTITIES_TYPES.entity) {
        return this.entityFilterHintsOrDefault;
      }

      return this.alarmFilterHintsOrDefault;
    },
  },

  watch: {
    advancedJsonWasChanged(value) {
      if (value) {
        this.resetFilterValidator();
      }
    },
  },
  async created() {
    if (this.required && this.$validator) {
      this.$validator.attach({
        name: 'filter',
        rules: 'required:true',
        getter: () => !checkIfGroupIsEmpty(this.form),
        context: () => this,
        vm: this,
      });
    }

    await this.fetchFilterHints();
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
        console.warn(err);

        this.$popups.error({ text: this.$t('errors.default') });
      }
    },

    updateJson(advancedJson) {
      try {
        this.$emit('input', filterToForm(advancedJson));

        this.advancedJson = advancedJson;
      } catch (err) {
        console.warn(err);

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
