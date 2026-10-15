<template>
  <c-layout column>
    <c-layout class="mb-2" align="center">
      <c-expand-btn
        v-if="expandable"
        :model-value="expanded"
        :loading="expanding"
        color="grey"
        @update:model-value="expandResponse"
      />
      <span class="ml-2">
        {{ $t(`alarm.timeline.steps.${step._t}`) }}
      </span>
      <span v-if="step.message"> : {{ step.message }} </span>
    </c-layout>
    <v-expand-transition>
      <v-card v-show="expanded">
        <v-card-text>
          <c-request-text-information v-if="response" :value="response" />
          <span v-else>
            {{ $t('common.noResponse') }}
          </span>
        </v-card-text>
      </v-card>
    </v-expand-transition>
  </c-layout>
</template>

<script>
export default {
  emits: ['expand'],
  props: {
    step: {
      type: Object,
      required: true,
    },
    expandable: {
      type: Boolean,
      default: false,
    },
    expanded: {
      type: Boolean,
      default: false,
    },
    expanding: {
      type: Boolean,
      default: false,
    },
    response: {
      type: String,
      required: false,
    },
  },
  setup(props, { emit }) {
    const expandResponse = value => emit('expand', value);

    return {
      expandResponse,
    };
  },
};
</script>
