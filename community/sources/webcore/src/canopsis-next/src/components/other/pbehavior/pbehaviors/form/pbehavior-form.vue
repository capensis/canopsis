<template>
  <v-tabs
    slider-color="primary"
    centered
    fixed-tabs
  >
    <v-tab :class="{ 'error--text': hasGeneralError }">
      {{ $t('common.general') }}
    </v-tab>
    <v-tab :class="{ 'error--text': hasPatternsError }">
      {{ $tc('common.pattern', 2) }}
    </v-tab>

    <v-tab-item eager>
      <v-layout
        class="py-3"
        column
      >
        <pbehavior-general-form
          ref="general"
          v-field="form"
          :no-enabled="noEnabled"
          :with-start-on-trigger="withStartOnTrigger"
        />
        <c-enabled-color-picker-field
          v-field="form.color"
          :label="$t('modals.createPbehavior.steps.color.label')"
        />
        <c-collapse-panel
          class="mb-2"
          :title="$t('recurrenceRule.title')"
        >
          <recurrence-rule-form
            v-field="form.rrule"
            :start="form.tstart"
          />
          <pbehavior-recurrence-rule-exceptions-field
            class="mt-2"
            v-field="form.exdates"
            :exceptions="form.exceptions"
            with-exdate-type
            @update:exceptions="updateExceptions"
          />
        </c-collapse-panel>
        <c-collapse-panel
          class="mt-2"
          v-if="!noComments"
          :title="$tc('common.comment', 2)"
        >
          <pbehavior-comments-field v-field="form.comments" />
        </c-collapse-panel>
      </v-layout>
    </v-tab-item>
    <v-tab-item eager>
      <v-layout
        class="py-3"
        justify-center
      >
        <v-flex xs12>
          <pbehavior-patterns-form
            ref="patterns"
            v-field="form.patterns"
          />
        </v-flex>
      </v-layout>
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { formMixin } from '@/mixins/form';

import RecurrenceRuleForm from '@/components/forms/recurrence-rule/recurrence-rule-form.vue';
import PbehaviorRecurrenceRuleExceptionsField from '@/components/other/pbehavior/exceptions/fields/pbehavior-recurrence-rule-exceptions-field.vue';

import PbehaviorCommentsField from '../fields/pbehavior-comments-field.vue';

import PbehaviorGeneralForm from './pbehavior-general-form.vue';
import PbehaviorPatternsForm from './pbehavior-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: {
    RecurrenceRuleForm,
    PbehaviorRecurrenceRuleExceptionsField,
    PbehaviorGeneralForm,
    PbehaviorCommentsField,
    PbehaviorPatternsForm,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    noPattern: {
      type: Boolean,
      default: false,
    },
    noEnabled: {
      type: Boolean,
      default: false,
    },
    noComments: {
      type: Boolean,
      default: false,
    },
    withStartOnTrigger: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      hasGeneralError: false,
      hasPatternsError: false,
    };
  },
  mounted() {
    this.$watch(() => this.$refs.general.hasAnyError, (value) => {
      this.hasGeneralError = value;
    });

    this.$watch(() => this.$refs.patterns.hasAnyError, (value) => {
      this.hasPatternsError = value;
    });
  },
  methods: {
    updateExceptions(exceptions) {
      this.updateField('exceptions', exceptions);
    },
  },
};
</script>
