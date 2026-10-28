<template>
  <c-collapse-panel :expanded="expanded">
    <template #header>
      <span class="job-http-data-collapse-panel__header text-white">
        <strong class="text-body-2">{{ titlePrefix }}:</strong>
        <span>{{ parsedRaw.startLine }}</span>
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
      <c-object-treeview
        v-if="tree && bodyJsonObject"
        v-model:opened="openedItems"
        :root-name="titlePrefix"
        :value="bodyJsonObject"
        value-max-width="720px"
        strong-title
      >
        <template #append="{ item }">
          <c-copy-btn
            v-if="copyable && item.path === titlePrefix"
            :value="parsedRaw.body"
            :tooltip="t('common.copyToClipboard')"
            small
            icon-small
            @success="showCopySuccessPopup"
            @error="showCopyErrorPopup"
          />
        </template>
      </c-object-treeview>
      <c-json-treeview v-else-if="bodyJsonObject" :json-object="bodyJsonObject" />
      <pre
        v-else
        class="text-break pre-wrap"
      >{{ parsedRaw.body }}
</pre>
    </div>
  </c-collapse-panel>
</template>

<script>
import { computed, ref } from 'vue';

import { maskSensitiveHttpHeaders, parseRawHttp } from '@/helpers/request';

import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';

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
    expanded: {
      type: Boolean,
      default: false,
    },
    tree: {
      type: Boolean,
      default: false,
    },
    copyable: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const popups = usePopups();
    const openedItems = ref([props.titlePrefix]);
    const parsedRaw = computed(() => {
      try {
        return maskSensitiveHttpHeaders(parseRawHttp(props.raw));
      } catch (error) {
        console.error(error);

        return {
          startLine: props.raw,
          headers: null,
          body: '',
        };
      }
    });
    const bodyJsonObject = computed(() => {
      try {
        return JSON.parse(parsedRaw.value.body);
      } catch {
        return null;
      }
    });
    const showCopySuccessPopup = () => popups.success({ text: t('popups.copySuccess') });
    const showCopyErrorPopup = () => popups.error({ text: t('popups.copyError') });

    return {
      t,
      openedItems,
      parsedRaw,
      bodyJsonObject,
      showCopySuccessPopup,
      showCopyErrorPopup,
    };
  },
};
</script>

<style scoped>
.job-http-data-collapse-panel__header {
  display: inline-flex;
  align-items: baseline;
  column-gap: 8px;
}
</style>
