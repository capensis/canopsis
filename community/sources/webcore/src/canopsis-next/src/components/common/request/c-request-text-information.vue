<template>
  <div class="d-flex flex-column align-stretch">
    <request-information-content :text="data.request" />
    <request-information-content v-if="data.response" :text="data.response" class="mt-4" />
  </div>
</template>

<script>
import { computed } from 'vue';

import RequestInformationContent from './partials/request-information-content.vue';

export default {
  components: { RequestInformationContent },

  props: {
    value: {
      type: String,
      required: true,
    },
  },

  setup(props) {
    const data = computed(() => {
      const [request, response] = props.value.split('\n\n');

      return {
        request,
        response,
      };
    });

    return {
      data,
    };
  },
};
</script>
