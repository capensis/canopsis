<template>
  <v-row class="request-url-field ga-2 flex-nowrap" align="start" no-gutters>
    <v-col class="request-url-field__column" cols="6">
      <c-select
        :model-value="method"
        :items="availableMethods"
        :label="methodLabel || $t('common.method')"
        :name="methodFieldName"
        :disabled="disabled"
        required
        @update:model-value="$emit('update:method', $event)"
      />
    </v-col>

    <v-col class="request-url-field__column" cols="6">
      <c-payload-text-field
        :value="url"
        :label="urlLabel || $t('common.url')"
        :name="urlFieldName"
        :variables="urlVariables"
        :disabled="disabled"
        required
        @update:value="$emit('update:url', $event)"
      >
        <template v-if="helpText" #append="">
          <c-help-icon
            :text="helpText"
            icon="help"
            color="grey darken-1"
            left
          />
        </template>
      </c-payload-text-field>
    </v-col>
  </v-row>
</template>

<script>
import { computed } from 'vue';

import { REQUEST_METHODS } from '@/constants';

export default {
  emits: ['update:method', 'update:url'],
  props: {
    method: {
      type: String,
      required: false,
    },
    url: {
      type: String,
      required: false,
    },
    methodLabel: {
      type: String,
      required: false,
    },
    urlLabel: {
      type: String,
      required: false,
    },
    helpText: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'request',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    urlVariables: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const availableMethods = Object.values(REQUEST_METHODS);
    const methodFieldName = computed(() => `${props.name}.method`);
    const urlFieldName = computed(() => `${props.name}.url`);

    return {
      availableMethods,
      methodFieldName,
      urlFieldName,
    };
  },
};
</script>

<style lang="scss" scoped>
.request-url-field {
  &__column {
    flex: 0 0 calc(50% - 4px);
    max-width: calc(50% - 4px);
    min-width: 0;
  }
}
</style>
