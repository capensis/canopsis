<template>
  <div class="request-information-content">
    <div
      v-for="(line, index) in request.commonInformationLines"
      :key="index"
      class="request-information-content__line"
    >
      <strong>{{ line.name }}</strong>
      <template v-if="line.value !== undefined">: {{ line.value }}</template>
    </div>

    <div v-if="request.body" class="request-information-content__body">
      <c-json-treeview v-if="request.isJsonBody" :json="request.body" />
      <template v-else>{{ request.body }}</template>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue';

import { isValidJsonData } from '@/helpers/json';

export default {
  props: {
    text: {
      type: String,
      required: true,
    },
  },

  setup(props) {
    const request = computed(() => {
      const [commonInformation = '', body] = props.text.split('\r\n\r\n');

      return {
        commonInformationLines: commonInformation.split('\r\n').map((line) => {
          const [name, value] = line.split(/:(.*)/s);

          return {
            name,
            value: value?.trim(),
          };
        }),
        isJsonBody: isValidJsonData(body),
        body,
      };
    });

    return {
      request,
    };
  },
};
</script>

<style lang="scss" scoped>
.request-information-content {
  display: flex;
  flex-direction: column;

  &__line {
    overflow-wrap: anywhere;
    white-space: pre-wrap;
  }

  &__body {
    margin-top: 20px;
    overflow-wrap: anywhere;
    white-space: pre-wrap;
  }
}
</style>
