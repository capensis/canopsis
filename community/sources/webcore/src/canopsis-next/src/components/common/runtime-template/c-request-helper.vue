<template>
  <c-runtime-template :template="template" />
</template>

<script>
import { get, unescape } from 'lodash';
import { computed, ref, onMounted, onBeforeUnmount } from 'vue';
import axios from 'axios';

import { RESPONSE_STATUSES } from '@/constants';

import { sanitizeHtml } from '@/helpers/html';

import { useI18n } from '@/hooks/i18n';
import { usePendingHandler } from '@/hooks/query/pending';

export default {
  props: {
    helperId: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    let cancelTokenSource = null;

    const { t } = useI18n();
    const template = ref('');
    const helperData = computed(() => window._handlebarsRequestHelper[props.helperId]);

    const { pending, handler: sendRequest } = usePendingHandler(async () => {
      if (!helperData.value) {
        return;
      }

      /**
       * Cancel previous request if exists
       */
      if (cancelTokenSource) {
        cancelTokenSource.cancel('Request cancelled due to new request');
      }

      /**
       * Create new cancel token
       */
      cancelTokenSource = axios.CancelToken.source();

      try {
        const {
          fn,
          hash: {
            method = 'get',
            url,
            headers,
            path,
            data,
            variable,
            username,
            password,
          },
        } = helperData.value.options;

        const axiosOptions = {
          method,
          url: unescape(url),
          cancelToken: cancelTokenSource.token,
        };

        if (headers) {
          axiosOptions.headers = JSON.parse(headers);
        }

        if (username || password) {
          axiosOptions.auth = { username, password };
        }

        if (data) {
          axiosOptions.data = JSON.parse(data);
        }

        const { data: responseData } = await axios(axiosOptions);

        const value = path ? get(responseData, path) : responseData;
        const context = variable ? { [variable]: value } : value;

        const html = await fn(context);

        /**
         * We need to use `all-inherit` class here because we may have multiple root elements in the response html
         */
        template.value = `<div class="all-inherit">${sanitizeHtml(html)}</div>`;
      } catch (err) {
        /**
         * Ignore cancelled requests
         */
        if (axios.isCancel(err)) {
          return;
        }

        console.error(err);

        const { status } = err.response || {};

        switch (status) {
          case RESPONSE_STATUSES.unauthorized:
            template.value = `<div class="all-inherit">${t('handlebars.requestHelper.errors.unauthorized')}</div>`;
            break;
          case RESPONSE_STATUSES.timeout:
            template.value = `<div class="all-inherit">${t('handlebars.requestHelper.errors.timeout')}</div>`;
            break;
          default:
            template.value = `<div class="all-inherit">${t('handlebars.requestHelper.errors.other')}</div>`;
        }
      } finally {
        cancelTokenSource = null;
      }
    });

    onMounted(sendRequest);

    onBeforeUnmount(() => {
      /**
       * Cancel ongoing request if exists
       */
      if (cancelTokenSource) {
        cancelTokenSource.cancel('Component unmounted');
        cancelTokenSource = null;
      }

      delete window._handlebarsRequestHelper[props.helperId];
    });

    return {
      pending,
      template,
    };
  },
};
</script>
