<template>
  <v-layout column>
    <v-layout
      class="mb-2"
      align-center
    >
      <c-expand-btn
        v-if="expandable"
        :expanded="expanded"
        :loading="expanding"
        color="grey"
        @expand="expandResponse"
      />
      <span class="ml-2">{{ $t(`alarm.timeLine.types.${step._t}`) }}</span>
      <span v-if="step.message">: {{ step.message }}</span>
    </v-layout>
    <v-expand-transition>
      <v-card v-show="expanded">
        <v-card-text>
          <c-request-text-information
            v-if="response"
            :value="response"
          />
          <span v-else>{{ $t('common.noResponse') }}</span>
        </v-card-text>
      </v-card>
    </v-expand-transition>
  </v-layout>
</template>

<script>
export default {
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
