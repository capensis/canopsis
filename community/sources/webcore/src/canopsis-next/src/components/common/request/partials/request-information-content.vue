<template>
  <div class="request-content">
    <request-information-content-row
      v-for="(line, index) in request.commonInformationLines"
      :key="index"
      :row="line"
    />
    <div
      v-if="request.body"
      class="request-content__body"
    >
      <json-treeview
        v-if="request.isJsonBody"
        :json="request.body"
      />
      <template v-else>
        {{ request.body }}
      </template>
    </div>
  </div>
</template>

<script>
import { isValidJsonData } from '@/helpers/json';

import JsonTreeview from '@/components/common/request/c-json-treeview.vue';

import RequestInformationContentRow from './request-information-content-row.vue';

export default {
  components: { RequestInformationContentRow, JsonTreeview },
  props: {
    text: {
      type: String,
      required: true,
    },
  },
  computed: {
    request() {
      const [commonInformation, body] = this.text.split('\r\n\r\n');

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
    },
  },
};
</script>

<style lang="scss">
.request-content {
  display: flex;
  flex-direction: column;

  &__body {
    margin-top: 20px;
    word-break: break-all;
  }
}
</style>
