<template>
  <v-layout
    class="llm-prompts-history-all-prompts-table"
    column
  >
    <v-layout
      class="pa-3 gap-3"
      align-center
      wrap
    >
      <v-flex
        xs12
        md5
      >
        <v-layout align-end>
          <v-text-field
            v-model="localSearchModel"
            :placeholder="$t('llm.promptsHistory.searchPlaceholder')"
            :aria-label="$t('llm.promptsHistory.searchPlaceholder')"
            class="mr-2"
            hide-details
            single-line
            dense
            @keydown.enter.prevent="emitSubmitSearch"
          />
          <c-action-btn
            :tooltip="$t('common.search')"
            :aria-label="$t('common.search')"
            icon="search"
            @click="emitSubmitSearch"
          />
          <c-action-btn
            :tooltip="$t('common.clearSearch')"
            :aria-label="$t('common.clearSearch')"
            icon="clear"
            @click="emitClearSearch"
          />
        </v-layout>
      </v-flex>
      <v-flex
        xs12
        md3
      >
        <v-switch
          :input-value="query.not_related_to_canopsis"
          :label="$t('llm.promptsHistory.notRelatedToCanopsis')"
          class="mt-0"
          color="primary"
          hide-details
          @change="onNotRelatedChange"
        />
      </v-flex>
      <v-flex
        xs12
        md3
      >
        <v-switch
          :input-value="query.group_by_chat"
          :label="$t('llm.promptsHistory.groupByChat')"
          class="mt-0"
          color="primary"
          hide-details
          @change="onGroupByChatChange"
        />
      </v-flex>
    </v-layout>

    <v-layout justify-center>
      <v-flex xs12>
        <v-card>
          <v-card-text class="pa-0">
            <c-advanced-data-table
              :headers="headers"
              :items="data"
              :loading="pending"
              :total-items="totalItems"
              :options="options"
              item-key="_id"
              hide-actions
              advanced-pagination
              @update:options="onUpdateOptions"
            >
              <template #datetime="{ item }">
                {{ item.created | date('long', '-') }}
              </template>
              <template #canopsis_related="{ item }">
                <c-enabled :value="item.canopsis_related" />
              </template>
              <template #prompt="{ item }">
                <span class="llm-prompts-history-all-prompts-table__prompt">{{ item.prompt }}</span>
              </template>
              <template #see_chat="{ item }">
                <c-action-btn
                  :tooltip="$t('llm.promptsHistory.seeChat')"
                  :aria-label="$t('llm.promptsHistory.seeChat')"
                  icon="forum"
                  @click="onSeeChat(item)"
                />
              </template>
            </c-advanced-data-table>
          </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

export default {
  props: {
    query: {
      type: Object,
      required: true,
    },
    localSearch: {
      type: String,
      default: '',
    },
    data: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      default: 0,
    },
    options: {
      type: Object,
      required: true,
    },
    headers: {
      type: Array,
      required: true,
    },
  },
  setup(props, { emit }) {
    const localSearchModel = computed({
      get: () => props.localSearch,
      set: (value) => {
        emit('update:localSearch', value);
      },
    });

    const emitSubmitSearch = () => {
      emit('submit-search');
    };

    const emitClearSearch = () => {
      emit('clear-search');
    };

    const onNotRelatedChange = (value) => {
      emit('not-related-change', value);
    };

    const onGroupByChatChange = (value) => {
      emit('group-by-chat-change', value);
    };

    const onUpdateOptions = (opts) => {
      emit('update:options', opts);
    };

    const onSeeChat = (item) => {
      emit('see-chat', item);
    };

    return {
      localSearchModel,
      emitSubmitSearch,
      emitClearSearch,
      onNotRelatedChange,
      onGroupByChatChange,
      onUpdateOptions,
      onSeeChat,
    };
  },
};
</script>

<style lang="scss">
.llm-prompts-history-all-prompts-table__prompt {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
