<template>
  <c-collapse-panel>
    <template #header>
      <span class="white--text">
        <strong>{{ titlePrefix }} : </strong>
        {{ parsedRaw.startLine }}
      </span>
    </template>
    <template v-if="parsedRaw.headers">
      <c-information-block-row
        v-for="(value, key) in parsedRaw.headers"
        :key="key"
        :label="key"
        :value="value"
      />
    </template>
    <div v-if="parsedRaw.body" class="mt-2">
      <c-json-treeview
        v-if="bodyJsonObject"
        :json-object="bodyJsonObject"
      />
      <pre v-else class="text-break pre-wrap">{{ parsedRaw.body }}</pre>
    </div>
  </c-collapse-panel>
</template>

<script>
import { computed } from 'vue';

import { parseRawHttp } from '@/helpers/request';

export default {
  props: {
    titlePrefix: {
      type: String,
      default: '',
    },
    raw: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const parsedRaw = computed(() => {
      try {
        return parseRawHttp(props.raw);
      } catch (error) {
        console.error(error);

        return error;
      }
    });

    const bodyJsonObject = computed(() => {
      try {
        return JSON.parse(parsedRaw.value.body);
      } catch (error) {
        return null;
      }
    });

    return {
      parsedRaw,
      bodyJsonObject,
    };
  },
};
</script>

<style lang="scss" scoped>
.ticket-status-jobs-http-data-collapse-panel__data {
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 0.875rem;
}
</style>
