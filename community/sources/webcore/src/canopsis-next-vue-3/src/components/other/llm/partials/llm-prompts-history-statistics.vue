<template>
  <v-row class="llm-prompts-history-statistics mt-4" justify="center" density="comfortable">
    <v-col cols="auto">
      <llm-prompts-history-statistic-big-block
        :title="t('llm.promptsHistory.statistics.totalInputTokens')"
        :value="statistics.input_tokens"
        :period="period"
        color="info"
      >
        <v-row class="ml-4">
          <v-col
            v-for="item in inputItems"
            :key="item.title"
            cols="auto"
          >
            <llm-prompts-history-statistic-small-block
              :title="item.title"
              :value="item.value"
              :percent="item.percent"
            />
          </v-col>
        </v-row>
      </llm-prompts-history-statistic-big-block>
    </v-col>

    <v-col cols="auto">
      <llm-prompts-history-statistic-big-block
        :title="t('llm.promptsHistory.statistics.totalOutputTokens')"
        :value="statistics.output_tokens"
        :period="period"
        color="warning"
      />
    </v-col>
  </v-row>
</template>

<script>
import { computed, ref, watch } from 'vue';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';
import { usePendingHandler } from '@/hooks/query/pending';
import { useLlm } from '@/hooks/store/modules/llm';

import LlmPromptsHistoryStatisticBigBlock from './llm-prompts-history-statistic-big-block.vue';
import LlmPromptsHistoryStatisticSmallBlock from './llm-prompts-history-statistic-small-block.vue';

const PERIOD_FORMAT = 'D MMM';

const DEFAULT_STATISTICS = {
  input_tokens: 0,
  output_tokens: 0,
  uncached_input_tokens: 0,
  cache_write_tokens: 0,
  cache_read_tokens: 0,
  uncached_input_tokens_ratio: 0,
  cache_write_tokens_ratio: 0,
  cache_read_tokens_ratio: 0,
};

export default {
  components: {
    LlmPromptsHistoryStatisticBigBlock,
    LlmPromptsHistoryStatisticSmallBlock,
  },

  props: {
    llmId: {
      type: String,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { fetchLlmStatsWithoutStore } = useLlm();
    const statistics = ref({ ...DEFAULT_STATISTICS });

    const { handler: fetchStatistics } = usePendingHandler(async () => {
      if (!props.llmId) {
        statistics.value = { ...DEFAULT_STATISTICS };

        return;
      }

      const response = await fetchLlmStatsWithoutStore({ id: props.llmId });
      const payload = response?.data && typeof response.data === 'object' && !Array.isArray(response.data)
        ? { ...response, ...response.data }
        : response;

      statistics.value = {
        ...DEFAULT_STATISTICS,
        ...payload,
      };
    });

    const period = computed(() => {
      if (!statistics.value.timestamp) {
        return '';
      }

      return `${convertDateToString(statistics.value.timestamp, PERIOD_FORMAT)} - ${convertDateToString(Date.now(), PERIOD_FORMAT)}`;
    });

    const inputItems = computed(() => [
      {
        title: t('llm.promptsHistory.statistics.uncachedTokens'),
        value: statistics.value.uncached_input_tokens,
        percent: statistics.value.uncached_input_tokens_ratio,
      },
      {
        title: t('llm.promptsHistory.statistics.cacheReadTokens'),
        value: statistics.value.cache_read_tokens,
        percent: statistics.value.cache_read_tokens_ratio,
      },
      {
        title: t('llm.promptsHistory.statistics.cacheWriteTokens'),
        value: statistics.value.cache_write_tokens,
        percent: statistics.value.cache_write_tokens_ratio,
      },
    ]);

    watch(() => props.llmId, fetchStatistics, { immediate: true });

    return {
      t,
      statistics,
      period,
      inputItems,
    };
  },
};
</script>
