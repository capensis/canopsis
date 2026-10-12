<template>
  <c-layout gap="2" column>
    <c-enabled-field
      v-if="!noEnabled"
      :model-value="form.enabled"
      hide-details
      with-background
      @update:model-value="updateField('enabled', $event)"
    />

    <c-form-general-patterns-tabs :form="form">
      <template #general="{ setRef }">
        <pbehavior-general-form
          :ref="setRef"
          :form="form"
          :no-enabled="noEnabled"
          :no-comments="noComments"
          :no-timezone="noTimezone"
          :with-start-on-trigger="withStartOnTrigger"
          :with-inherited="withInherited"
          :name-label="nameLabel"
          :name-tooltip="nameTooltip"
          @update:form="updateModel"
        />
      </template>

      <template #patterns="{ setRef }">
        <pbehavior-patterns-form
          :ref="setRef"
          :form="form.patterns"
          :pbehavior-id="pbehaviorId"
          :pbehavior-counter-type="pbehaviorCounterType"
          @update:form="updateField('patterns', $event)"
        />
      </template>
    </c-form-general-patterns-tabs>
  </c-layout>
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';

import PbehaviorGeneralForm from './pbehavior-general-form.vue';
import PbehaviorPatternsForm from './pbehavior-patterns-form.vue';

export default {
  components: {
    PbehaviorGeneralForm,
    PbehaviorPatternsForm,
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

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateModel, updateField } = useModelField(props, emit);

    return {
      updateModel,
      updateField,
    };
  },
};
</script>
