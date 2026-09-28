<template>
  <v-layout class="gap-2" column>
    <c-enabled-field
      v-if="!noEnabled"
      v-field="form.enabled"
      hide-details
      with-background
    />
    <c-form-general-patterns-tabs :form="form">
      <template #general="{ setRef }">
        <pbehavior-general-form
          v-field="form"
          :ref="setRef"
          :no-enabled="noEnabled"
          :no-comments="noComments"
          :no-timezone="noTimezone"
          :with-start-on-trigger="withStartOnTrigger"
          :with-inherited="withInherited"
          :name-label="nameLabel"
          :name-tooltip="nameTooltip"
        />
      </template>
      <template #patterns="{ setRef }">
        <pbehavior-patterns-form
          v-field="form.patterns"
          :ref="setRef"
          :pbehavior-id="pbehaviorId"
          :pbehavior-counter-type="pbehaviorCounterType"
        />
      </template>
    </c-form-general-patterns-tabs>
  </v-layout>
</template>

<script>
import PbehaviorGeneralForm from './pbehavior-general-form.vue';
import PbehaviorPatternsForm from './pbehavior-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: {
    PbehaviorGeneralForm,
    PbehaviorPatternsForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    pbehaviorId: {
      type: String,
      required: false,
    },
    form: {
      type: Object,
      required: true,
    },
    noEnabled: {
      type: Boolean,
      default: false,
    },
    noComments: {
      type: Boolean,
      default: false,
    },
    noTimezone: {
      type: Boolean,
      default: false,
    },
    withStartOnTrigger: {
      type: Boolean,
      default: false,
    },
    withInherited: {
      type: Boolean,
      default: false,
    },
    nameLabel: {
      type: String,
      required: false,
    },
    nameTooltip: {
      type: String,
      required: false,
    },
    pbehaviorCounterType: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
