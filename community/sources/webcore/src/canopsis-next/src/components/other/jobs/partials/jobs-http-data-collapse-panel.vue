<template>
  <c-collapse-panel>
    <template #header>
      <span class="white--text">
        <strong>{{ titlePrefix }} : </strong>
        {{ parsedRaw.title }}
      </span>
    </template>
    <template v-if="parsedRaw.headers.length">
      <c-information-block-row
        v-for="(row, index) in parsedRaw.headers"
        :key="index"
        :label="row.key"
        :value="row.value"
      />
    </template>
    <div v-if="parsedRaw.data">
      <c-json-treeview
        v-if="dataJsonObject"
        :json-object="dataJsonObject"
      />
      <pre v-else class="text-break pre-wrap">{{ parsedRaw.data }}</pre>
    </div>
    <template v-else>
      <span class="grey--text">{{ $t('common.noDataAvailable') }}</span>
    </template>
  </c-collapse-panel>
</template>

<script>
import { computed } from 'vue';

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
      const [headersRaw = '', data = ''] = props.raw.split('\n\n');
      const headers = headersRaw.split('\n');

      return {
        title: headers.shift(),
        headers: headers.map((header) => {
          const [key, value] = header.split(/\s*:\s*/);

          return { key, value };
        }),
        data,
      };
    });

    const dataJsonObject = computed(() => {
      try {
        return JSON.parse(parsedRaw.value.data);
      } catch (error) {
        return null;
      }
    });

    return {
      parsedRaw,
      dataJsonObject,
    };
  },
};
</script>

<style lang="scss" scoped>
.jobs-http-data-collapse-panel__data {
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 0.875rem;
}
</style>
